package controller

import (
	"errors"
	"net/http"

	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) UpvotePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	err = ctr.store.CreateUpvote(c.Context(), id, user.ID)
	if errors.Is(err, store.ErrNoRowsAffected) {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Already upvoted").Send(c)
	}
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "Post with such id was not found").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) RevertPostUpvote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.DeleteUpvote(c.Context(), id, user.ID); err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) DownvotePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
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
		return httperr.New(codes.Omit, http.StatusNotFound, "Post with such id was not found").Send(c)
	}
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) RevertPostDownvote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Invalid `id` param").Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if err = ctr.store.DeleteDownvote(c.Context(), id, user.ID); err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
