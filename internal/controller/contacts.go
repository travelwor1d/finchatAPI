package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func (ctr *Ctr) ListContacts(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `size` param").Send(c)
	}
	contacts, err := ctr.store.ListContacts(c.Context(), user.ID, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if contacts == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"contacts": []interface{}{}})
	}
	return c.JSON(fiber.Map{"contacts": contacts})
}

func (ctr *Ctr) GetContact(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
	}
	contact, err := ctr.store.GetContact(c.Context(), user.ID, id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "Contact with such id was not found").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(contact)
}

type createContactPayload struct {
	ContactID int `json:"contactId" validate:"required"`
}

func (ctr *Ctr) CreateContact(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	var p createContactPayload
	if err := c.BodyParser(&p); err != nil {
		return errParseBody.SetDetail(err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}
	contact, err := ctr.store.CreateContact(c.Context(), user.ID, p.ContactID)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "No user with such id").Send(c)
	}
	if errors.Is(err, store.ErrAlreadyExists) {
		return httperr.New(codes.Omit, http.StatusConflict, "This user is already in your contacts").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(contact)
}

func (ctr *Ctr) DeleteContact(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
	}

	err = ctr.store.DeleteContact(c.Context(), user.ID, id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "Contact with such id was not found").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
