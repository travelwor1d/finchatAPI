package controller

import (
	"errors"
	"net/http"

	"github.com/finchatapp/finchat-api/internal/logerr"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) UpvotePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	err = ctr.store.CreateUpvote(c.Context(), id, user.ID)
	if errors.Is(err, store.ErrNoRowsAffected) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "already upvoted").Send(c)
	}
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "post not found").Send(c)
	}
	if err != nil {
		logerr.LogError(err)
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) RevertPostUpvote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.DeleteUpvote(c.Context(), id, user.ID); err != nil {
		logerr.LogError(err)
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) DownvotePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	err = ctr.store.CreateDownvote(c.Context(), id, user.ID)
	if errors.Is(err, store.ErrNoRowsAffected) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "already downvoted").Send(c)
	}
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "post not found").Send(c)
	}
	if err != nil {
		logerr.LogError(err)
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) RevertPostDownvote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.DeleteDownvote(c.Context(), id, user.ID); err != nil {
		logerr.LogError(err)
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
