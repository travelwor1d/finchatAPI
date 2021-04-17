package controller

import (
	"errors"
	"fmt"
	"net/http"

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
	creds, err := ctr.store.GetUserCredsByEmail(c.Context(), p.Email)
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
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"required|email"`
	Password  string `json:"password" validate:"required|minLen:6"`
}

func (ctr *Ctr) Register(c *fiber.Ctx) error {
	var p registerPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return httperr.New(codes.Omit, http.StatusUnprocessableEntity, v.Errors.One()).Send(c)
	}

	var userType string
	inviteCode := c.Query("inviteCode")
	if inviteCode != "" {
		userType = "GOAT"
		if len(inviteCode) != 6 {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Invite code is 6 chars long string").Send(c)
		}
		status, found, err := ctr.store.GetInviteCodeStatus(c.Context(), inviteCode)
		if err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
		if !found {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid invite code").Send(c)
		}
		if status == "USED" {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Invite code was already used").Send(c)
		}
		if status == "EXPIRED" {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Invite code expired")
		}
	} else {
		userType = "USER"
	}

	user := &model.User{FirstName: p.FirstName, LastName: p.LastName, Phone: p.Phone, Email: p.Email, Type: userType}
	user, err := ctr.store.CreateUser(c.Context(), user, p.Password, inviteCode)
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
	token, err := ctr.jwtManager.Generate(fmt.Sprint(user.ID))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"user": user, "token": token})
}

func matches(hash, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return false
	}
	return true
}
