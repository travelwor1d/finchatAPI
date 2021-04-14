package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) UpvotePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.CreateUpvote(c.Context(), id, userID); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) RevertPostUpvote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.DeleteUpvote(c.Context(), id, userID); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) DownvotePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.CreateDownvote(c.Context(), id, userID); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) RevertPostDownvote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	userID, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.DeleteDownvote(c.Context(), id, userID); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
