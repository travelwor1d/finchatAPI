package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func (ctr *Ctr) ListThreads(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	threads, err := ctr.store.ListThreads(c.Context(), user.ID, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if threads == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"threads": []interface{}{}})
	}
	return c.JSON(fiber.Map{"threads": threads})

}

func (ctr *Ctr) GetThread(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	thread, err := ctr.store.GetThread(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "no thread with such id").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(thread)
}

type createThreadPayload struct {
	Title        string `json:"title" validate:"-"`
	Type         string `json:"type" validate:"required"`
	Participants []int  `json:"participants" validate:"required|ints"`
}

func (ctr *Ctr) CreateThread(c *fiber.Ctx) error {
	var p createThreadPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	thread, err := ctr.store.CreateThread(c.Context(), &model.Thread{
		Title: p.Title,
		Type:  p.Type,
	}, append(p.Participants, user.ID))
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(thread)
}

type sendMessagePayload struct {
	Message string  `json:"message" validate:"required"`
	Type    *string `json:"type" validate:"-"`
}

func (ctr *Ctr) SendMessage(c *fiber.Ctx) error {
	var p sendMessagePayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	threadID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}

	timestamp, err := ctr.msg.User(user.ID).SendMessage(c.Context(), threadID, user.ID, p.Message)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{
		"message":   p.Message,
		"timestamp": timestamp,
	})
}
