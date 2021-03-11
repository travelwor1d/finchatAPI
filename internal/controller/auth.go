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

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ctr *Ctr) Login(c *fiber.Ctx) error {
	var p LoginPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	creds, err := ctr.store.GetUserCredsByEmail(c.Context(), p.Email)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.InvalidCredentials, http.StatusNotFound, "user not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if !matches([]byte(creds.Hash), []byte(p.Password)) {
		return httperr.New(codes.InvalidCredentials, http.StatusBadRequest, "passwords did not match").Send(c)
	}
	token, err := ctr.jwtManager.Generate(fmt.Sprint(creds.UserID))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"token": token})
}

type RegisterPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"required|email"`
	Password  string `json:"password" validate:"required|minLen:6"`
}

func (ctr *Ctr) Register(c *fiber.Ctx) error {
	var p RegisterPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	var userType string
	inviteCode := c.Query("inviteCode")
	if inviteCode != "" {
		userType = "GOAT"
		if len(inviteCode) != 6 {
			return httperr.New(codes.Omit, http.StatusBadRequest, "invite code is 6 chars long string").Send(c)
		}
		status, found, err := ctr.store.GetInviteCodeStatus(c.Context(), inviteCode)
		if err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
		if !found {
			return httperr.New(codes.Omit, http.StatusBadRequest, "invalid invite code").Send(c)
		}
		if status == "USED" {
			return httperr.New(codes.Omit, http.StatusBadRequest, "invite code was already used").Send(c)
		}
		if status == "EXPIRED" {
			return httperr.New(codes.Omit, http.StatusBadRequest, "invite code expired")
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
			"user with provided email already exists",
		).Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	token, err := ctr.jwtManager.Generate(fmt.Sprint(user.ID))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"token": token, "verified": user.Verified})
}

func matches(hash, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return false
	}
	return true
}
