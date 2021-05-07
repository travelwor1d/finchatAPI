package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func (ctr *Ctr) ListUsers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `size` param").Send(c)
	}
	q := c.Query("query")
	userTypes := c.Query("userTypes")
	userTypes, err = getUserTypes(userTypes)
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, err.Error()).Send(c)
	}
	users, err := ctr.store.SearchUsers(c.Context(), q, userTypes, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if users == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"users": []interface{}{}})
	}
	return c.JSON(fiber.Map{"users": users})
}

func (ctr *Ctr) GetUser(c *fiber.Ctx) error {
	var id int
	var user *model.User
	var err error
	var httpErr *httperr.HTTPErr
	if strings.HasSuffix(c.Path(), "/me") {
		user, httpErr = ctr.userFromCtx(c)
		if err != nil {
			return httpErr.Send(c)
		}
	} else {
		id, err = strconv.Atoi(c.Params("id"))
		if err != nil {
			return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
		}
		user, err = ctr.store.GetUser(c.Context(), id)
		if errors.Is(err, store.ErrNotFound) {
			return httperr.New(codes.Omit, http.StatusNotFound, "User was not found. Please try again").Send(c)
		}
		if err != nil {
			ctr.lr.LogError(err, c.Request())
			return errInternal.SetDetail(err).Send(c)
		}
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
		return errParseBody.SetDetail(err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.UpdateUser(c.Context(), user.ID, p.FirstName, p.LastName, p.ProfileAvatar)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "User was not found").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(user)
}

func (ctr *Ctr) SoftDeleteUser(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	err := ctr.store.SoftDeleteUser(c.Context(), user.ID)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) UndeleteUser(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	err := ctr.store.UndeleteUser(c.Context(), user.ID)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "User has not been deleted").Send(c)
	}
	if errors.Is(err, store.ErrUserNotDeleted) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "User has not been deleted").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
