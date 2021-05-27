package controller

import (
	"errors"
	"net/http"
	"strconv"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
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

	req := srequest.NewGridList{
		PageSize:   size,
		PageNumber: page,
	}
	req.CustomFilters = make([]srequest.CustomFilterItem, 0)
	req.CustomFilters = append(req.CustomFilters, srequest.CustomFilterItem{
		Name:   "onlyContacts",
		Values: []string{"true"},
	})
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
	paging, contacts, err := ctr.user.SearchUsers(c.Context(), user.ID, req)

	if contacts == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"contacts": []interface{}{}, "paging": new(interface{})})
	}
	return c.JSON(fiber.Map{"contacts": contacts, "paging": paging})
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

	req := srequest.NewGridList{
		PageSize:   1,
		PageNumber: 1,
	}
	req.CustomFilters = make([]srequest.CustomFilterItem, 0)
	req.CustomFilters = append(req.CustomFilters, srequest.CustomFilterItem{
		Name:   "contactId",
		Values: []string{strconv.Itoa(id)},
	})
	req.Sorts = make([]srequest.SortItem, 0)
	req.Sorts = append(req.Sorts, srequest.SortItem{
		Field: "last_name",
		Dir:   "asc",
	})
	_, contacts, err := ctr.user.SearchUsers(c.Context(), user.ID, req)
	if errors.Is(err, store.ErrNotFound) || len(contacts) == 0 {
		return httperr.New(codes.Omit, http.StatusNotFound, "Contact with such id was not found").Send(c)
	}

	contact := contacts[0]
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
	contact, err := ctr.contact.CreateContact(c.Context(), user.ID, p.ContactID)
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

	err = ctr.contact.DeleteContact(c.Context(), user.ID, id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "Contact with such id was not found").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
