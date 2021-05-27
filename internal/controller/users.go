package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
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
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}

	req := srequest.NewGridList{
		PageSize:   size,
		PageNumber: page,
	}
	req.CustomFilters = make([]srequest.CustomFilterItem, 0)
	req.CustomFilters = append(req.CustomFilters, srequest.CustomFilterItem{
		Name:   "userTypes",
		Values: getUserTypes(c.Query("userTypes")),
	})
	{
		customFilters := strings.Split(c.Query("filters", ""), ",")
		for _, v := range customFilters {
			if v != "" {
				req.CustomFilters = append(req.CustomFilters, srequest.CustomFilterItem{
					Name:   v,
					Values: []string{"true"},
				})
			}
		}
	}
	req.CustomFilters = append(req.CustomFilters, srequest.CustomFilterItem{
		Name:   "query",
		Values: []string{c.Query("query")},
	})
	req.Sorts = make([]srequest.SortItem, 0)
	req.Sorts = append(req.Sorts, srequest.SortItem{
		Field: "last_name",
		Dir:   "asc",
	})
	req.Sorts = append(req.Sorts, srequest.SortItem{
		Field: "first_name",
		Dir:   "asc",
	})
	paging, users, err := ctr.user.SearchUsers(c.Context(), user.ID, req)

	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if users == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"users": []interface{}{}, "paging": new(interface{})})
	}
	return c.JSON(fiber.Map{"users": users, "paging": paging})
}

func (ctr *Ctr) GetUser(c *fiber.Ctx) error {
	var id int
	var user *model.User
	var err error
	var httpErr *httperr.HTTPErr
	if strings.HasSuffix(c.Path(), "/me") {
		user, httpErr = ctr.userFromCtx(c)
		if httpErr != nil {
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
	Username      *string `json:"username" validate:"-"`
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
	user, err := ctr.store.UpdateUser(c.Context(), user.ID, p.FirstName, p.LastName, p.Username, p.ProfileAvatar)
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
