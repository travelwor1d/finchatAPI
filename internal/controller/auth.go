package controller

import (
	"errors"
	"net/http"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

type registerPayload struct {
	FirstName string `json:"firstName" validate:"required|alpha"`
	LastName  string `json:"lastName" validate:"required|alpha"`
	Phone
	Email string `json:"email" validate:"required|email"`
}

func (ctr *Ctr) Register(c *fiber.Ctx) error {
	var p registerPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
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

	user := &model.User{
		FirstName: p.FirstName, LastName: p.LastName, Phonenumber: p.formattedPhonenumber(), CountryCode: p.CountryCode, Email: p.Email, Type: userType,
	}
	user, err := ctr.store.CreateUser(c.Context(), user, inviteCode)
	if errors.Is(err, store.ErrAlreadyExists) {
		return httperr.New(
			codes.EmailAlreadyTaken,
			http.StatusBadRequest,
			"User with provided email already exists",
		).Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if err := ctr.msg.User(user.ID).Register(c.Context(), user.FirstName, user.LastName, user.Email); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"user": user})
}

func (ctr *Ctr) EmailValidation(c *fiber.Ctx) error {
	email := c.Query("email")
	if !validate.IsEmail(email) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Please enter a valid email").Send(c)
	}
	taken, err := ctr.store.IsEmailTaken(c.Context(), email)
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
