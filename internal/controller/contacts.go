package controller

import (
	"net/http"
	"strconv"

	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) ListContacts(c *fiber.Ctx) error {
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	if id != userID {
		return httperr.New(codes.Omit, http.StatusForbidden, "ids in path and token does not match").Send(c)
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
	return c.JSON(fiber.Map{"contacts": contacts})
}

func (ctr *Ctr) ListContactRequests(c *fiber.Ctx) error {
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	if id != userID {
		return httperr.New(codes.Omit, http.StatusForbidden, "ids in path and token does not match").Send(c)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	contacts, err := ctr.store.ListContactRequests(c.Context(), id, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"contactRequests": contacts})
}

type createContactRequest struct {
	ContactOwnerID int `json:"contactOwnerId" validate:"required"`
	ContactID      int `json:"contactId" validate:"required"`
}

func (ctr *Ctr) CreateContactRequest(c *fiber.Ctx) error {
	panic("unimplemented")
}

type patchContactRequest struct {
	Status string `json:"status" validate:"required"`
}

func (ctr *Ctr) PatchContactRequest(c *fiber.Ctx) error {
	panic("unimplemented")
}
