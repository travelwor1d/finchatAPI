package controller

import (
	"net/http"
	"strings"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

type registerPayload struct {
	FirstName string `json:"firstName" validate:"required|alpha"`
	LastName  string `json:"lastName" validate:"required|alpha"`
	Phone
	Email    string  `json:"email" validate:"required|email"`
	Username *string `json:"username" validate:"maxLength:15"`
}

func (ctr *Ctr) Register(c *fiber.Ctx) error {
	var p registerPayload
	if err := c.BodyParser(&p); err != nil {
		return errParseBody.SetDetail(err).Send(c)
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
			ctr.lr.LogError(err, c.Request())
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

	p.Email = sanitizeEmail(p.Email)

	user := &model.User{
		FirstName: p.FirstName, LastName: p.LastName,
		Phonenumber: p.formattedPhonenumber(), CountryCode: p.CountryCode,
		Email: p.Email, Username: p.Username,
		Type: userType,
	}
	isTaken, err := ctr.store.IsEmailTaken(c.Context(), p.Email)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if isTaken {
		return httperr.New(
			codes.Omit,
			http.StatusBadRequest,
			"User with provided email already exists",
		).Send(c)
	}
	isTaken, err = ctr.store.IsPhoneNumberTaken(c.Context(), p.formattedPhonenumber())
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if isTaken {
		return httperr.New(
			codes.Omit,
			http.StatusBadRequest,
			"User with provided phone number already exists",
		).Send(c)
	}
	if user.Username != nil {
		isTaken, err = ctr.store.IsUsernameTaken(c.Context(), *p.Username)
		if err != nil {
			ctr.lr.LogError(err, c.Request())
			return errInternal.SetDetail(err).Send(c)
		}
		if isTaken {
			return httperr.New(
				codes.Omit,
				http.StatusBadRequest,
				"User with provided username already exists",
			).Send(c)
		}
	}

	user, err = ctr.store.UpsertUser(c.Context(), user, inviteCode)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if err := ctr.msg.User(user.ID).Register(c.Context(), user.FirstName, user.LastName, user.Email); err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"user": user})
}

func (ctr *Ctr) EmailValidation(c *fiber.Ctx) error {
	email := c.Query("email")
	if !validate.IsEmail(email) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Please enter a valid email").Send(c)
	}
	taken, err := ctr.store.IsEmailTaken(c.Context(), sanitizeEmail(email))
	if err != nil {
		ctr.lr.LogError(err, c.Request())
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
	c.Request()
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if taken {
		return c.JSON(fiber.Map{"isTaken": true, "message": "A user already exists with this phone number"})
	}
	return c.JSON(fiber.Map{"isTaken": false, "message": ""})
}

func (ctr *Ctr) UsernameValidation(c *fiber.Ctx) error {
	username := c.Query("username")
	if !validate.IsAlphaNum(username) || !validate.MaxLength(username, 15) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Please enter a valid username").Send(c)
	}
	taken, err := ctr.store.IsUsernameTaken(c.Context(), username)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if taken {
		return c.JSON(fiber.Map{"isTaken": true, "message": "A user already exists with this username"})
	}
	return c.JSON(fiber.Map{"isTaken": false, "message": ""})
}

func sanitizeEmail(s string) string {
	return strings.ToLower(s)
}
