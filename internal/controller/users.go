package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) ListUsers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	q := c.Query("query")
	userTypes := c.Query("userTypes")
	userTypes, err = getUserTypes(userTypes)
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, err.Error()).Send(c)
	}
	users, err := ctr.store.SearchUsers(c.Context(), q, userTypes, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"users": users})
}

func (ctr *Ctr) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "user with such id was not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(user)
}

type updateUserPayload struct {
	FirstName     *string `json:"firstName" validate:"-"`
	LastName      *string `json:"lastName" validate:"-"`
	ProfileAvatar *string `json:"profileAvatar" validate:"-"`
}

func (ctr *Ctr) UpdateUser(c *fiber.Ctx) error {
	var p updateUserPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	claimsID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if id != claimsID {
		return httperr.New(codes.Omit, http.StatusForbidden, "ids in claims and path params did not match").Send(c)
	}
	user, err := ctr.store.UpdateUser(c.Context(), id, p.FirstName, p.LastName, p.ProfileAvatar)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "user with such id was not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(user)
}

func (ctr *Ctr) SoftDeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	claimsID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if id != claimsID {
		return httperr.New(codes.Omit, http.StatusForbidden, "ids in claims and path params did not match").Send(c)
	}
	err = ctr.store.SoftDeleteUser(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "user with such id was not found").Send(c)
	}
	if errors.Is(err, store.ErrUserDeleted) {
		return httperr.New(codes.Omit, http.StatusNotFound, "user was already deleted").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"success": true})
}

func (ctr *Ctr) UndeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	claimsID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if id != claimsID {
		return httperr.New(codes.Omit, http.StatusForbidden, "ids in claims and path params did not match").Send(c)
	}
	err = ctr.store.UndeleteUser(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "user with such id was not found").Send(c)
	}
	if errors.Is(err, store.ErrUserNotDeleted) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "user with such id has not been deleted").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"success": true})
}
