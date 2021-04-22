package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"golang.org/x/crypto/bcrypt"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ctr *Ctr) Login(c *fiber.Ctx) error {
	var p loginPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	creds, err := ctr.store.GetUserCredsByEmail(c.Context(), sanitizeEmail(p.Email))
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.InvalidCredentials, http.StatusNotFound, "The email you entered does not match an account").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if !matches([]byte(creds.Hash), []byte(p.Password)) {
		return httperr.New(codes.InvalidCredentials, http.StatusBadRequest, "The password you entered is incorrect").Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), creds.UserID)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	token, err := ctr.jwtManager.Generate(fmt.Sprint(creds.UserID))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"user": user, "token": token})
}

type registerPayload struct {
	FirstName string `json:"firstName" validate:"required|alpha"`
	LastName  string `json:"lastName" validate:"required|alpha"`
	Phone
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required"`
}

func (ctr *Ctr) Register(c *fiber.Ctx) error {
	var p registerPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	if !validate.MinLength(p.Password, 8) {
		return httperr.NewValidationErr(nil, "Your password should be at least 8 characters").Send(c)
	}

	var userType string
	inviteCode := c.Query("inviteCode")
	if inviteCode != "" {
		userType = "GOAT"
		if len(inviteCode) != 6 {
			return httperr.NewValidationErr(nil, "Your invitation code should be 6 chars long string").Send(c)
		}
		status, found, err := ctr.store.GetInviteCodeStatus(c.Context(), inviteCode)
		if err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
		if !found {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Your invitation code is invalid").Send(c)
		}
		if status == "USED" {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Your invitation code was already used").Send(c)
		}
		if status == "EXPIRED" {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Your invitation code has expired").Send(c)
		}
	} else {
		userType = "USER"
	}

	if v := validate.Struct(p); !v.Validate() {
		return httperr.NewValidationErr(v, "Please enter valid input data").Send(c)
	}
	if !p.Phone.Validate() {
		return httperr.NewValidationErr(nil, "Please enter a valid phone number").Send(c)
	}
	if !validate.MaxLength(p.FirstName, 50) || !validate.MaxLength(p.LastName, 50) {
		return httperr.NewValidationErr(nil, "First or last names on Finchat can't have too many characters").Send(c)
	}

	user := &model.User{
		FirstName: p.FirstName, LastName: p.LastName, Phonenumber: p.formattedPhonenumber(), CountryCode: p.CountryCode, Email: sanitizeEmail(p.Email), Type: userType,
	}
	user, err := ctr.store.CreateUser(c.Context(), user, p.Password, inviteCode)
	if errors.Is(err, store.ErrAlreadyExists) {
		return httperr.New(
			codes.EmailAlreadyTaken,
			http.StatusBadRequest,
			"User with provided email or phone number already exists",
		).Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if err := ctr.msg.User(user.ID).Register(c.Context(), user.FirstName, user.LastName, user.Email); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	token, err := ctr.jwtManager.Generate(fmt.Sprint(user.ID))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"user": user, "token": token})
}

func (ctr *Ctr) EmailValidation(c *fiber.Ctx) error {
	email := c.Query("email")
	if !validate.IsEmail(email) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Please enter a valid email").Send(c)
	}
	taken, err := ctr.store.IsEmailTaken(c.Context(), sanitizeEmail(email))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if taken {
		return c.JSON(fiber.Map{"isTaken": true, "message": "A user already exists with this email"})
	}
	return c.JSON(fiber.Map{"isTaken": false, "message": ""})
}

func (ctr *Ctr) PhonenumberValidation(c *fiber.Ctx) error {
	var q Phone
	if err := c.QueryParser(&q); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse query string").Send(c)
	}
	if v := validate.Struct(q); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}
	if !q.Validate() {
		return httperr.NewValidationErr(nil, "Please enter a valid phone number").Send(c)
	}

	taken, err := ctr.store.IsPhoneNumberTaken(c.Context(), q.formattedPhonenumber())
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if taken {
		return c.JSON(fiber.Map{"isTaken": true, "message": "A user already exists with this phone number"})
	}
	return c.JSON(fiber.Map{"isTaken": false, "message": ""})
}

func matches(hash, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return false
	}
	return true
}

func sanitizeEmail(s string) string {
	return strings.ToLower(s)
}
