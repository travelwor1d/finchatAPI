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

func (ctr *Ctr) ListContactsOrRequests(c *fiber.Ctx) error {
	requests := c.Query("requests", "false")
	if requests == "true" {
		return ctr.ListContactRequests(c)
	} else if requests == "false" {
		return ctr.ListContacts(c)
	}
	return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `requests` param").Send(c)
}

func (ctr *Ctr) ListContacts(c *fiber.Ctx) error {
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	contacts, err := ctr.store.ListContacts(c.Context(), id, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if contacts == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"contacts": []interface{}{}})
	}
	return c.JSON(fiber.Map{"contacts": contacts})
}

func (ctr *Ctr) GetContact(c *fiber.Ctx) error {
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	contact, err := ctr.store.GetContact(c.Context(), userID, id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(contact)
}

func (ctr *Ctr) ListContactRequests(c *fiber.Ctx) error {
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	contactRequests, err := ctr.store.ListContactRequests(c.Context(), id, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if contactRequests == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"contactRequests": []interface{}{}})
	}
	return c.JSON(fiber.Map{"contactRequests": contactRequests})
}

type createContactRequest struct {
	ContactID int `json:"contactId" validate:"required"`
}

func (ctr *Ctr) CreateContactRequest(c *fiber.Ctx) error {
	var p createContactRequest
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if p.ContactID == id {
		return httperr.New(codes.Omit, http.StatusBadRequest, "user cannot request a contact with themself").Send(c)
	}
	r, err := ctr.store.CreateContactRequest(c.Context(), id, p.ContactID)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(r)
}

type patchContactRequest struct {
	Status string `json:"status" validate:"required"`
}

func (ctr *Ctr) PatchContactRequest(c *fiber.Ctx) error {
	var p patchContactRequest
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}

	r, err := ctr.store.GetContactRequest(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "not found").Send(c)
	}
	if userID != r.ContactID {
		return httperr.New(codes.Omit, http.StatusForbidden, "cannot approve or deny not theirs contact request").Send(c)
	}

	if p.Status == "ACCEPTED" {
		if err := ctr.store.ApproveContactRequest(c.Context(), id); err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
	} else if p.Status == "DENIED" {
		if err := ctr.store.DenyContactRequest(c.Context(), id); err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
	} else {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid status").Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) DeleteContact(c *fiber.Ctx) error {
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}

	err = ctr.store.DeleteContact(c.Context(), userID, id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
