package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func (ctr *Ctr) ListPosts(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	posts, err := ctr.store.ListPosts(c.Context(), &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if posts == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"posts": []interface{}{}})
	}
	return c.JSON(fiber.Map{"posts": posts})
}

func (ctr *Ctr) GetPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	post, err := ctr.store.GetPost(c.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "post with such id was not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(post)
}

type createPostPayload struct {
	Title       string     `json:"title" validate:"required"`
	Content     string     `json:"content" validate:"required"`
	ImageURLs   []string   `json:"imageUrls" validate:"strings"`
	Tickers     []string   `json:"tickers" validate:"strings"`
	PublishedAt *time.Time `json:"publishedAt" validate:"-"`
}

func (ctr *Ctr) CreatePost(c *fiber.Ctx) error {
	var p createPostPayload
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
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}

	post, err := ctr.store.CreatePost(c.Context(), &model.Post{
		Title:       p.Title,
		Content:     p.Content,
		ImageURLs:   p.ImageURLs,
		Tickers:     p.Tickers,
		PostedBy:    user.ID,
		PublishedAt: p.PublishedAt,
	})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(post)
}

func (ctr *Ctr) ListComments(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `page` param").Send(c)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `size` param").Send(c)
	}
	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `postId` param").Send(c)
	}
	comments, err := ctr.store.ListComments(c.Context(), postID, &store.Pagination{Limit: size, Offset: size * (page - 1)})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if comments == nil {
		// Return an empty array.
		return c.JSON(fiber.Map{"comments": []interface{}{}})
	}
	return c.JSON(fiber.Map{"comments": comments})
}

func (ctr *Ctr) GetComment(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `postId` param").Send(c)
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `id` param").Send(c)
	}
	comment, err := ctr.store.GetComment(c.Context(), postID, id)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "comment with such id was not found").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(comment)
}

type createCommentPayload struct {
	Content string `json:"content" validate:"required"`
}

func (ctr *Ctr) CreateComment(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid `postId` param").Send(c)
	}
	var p createCommentPayload
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
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}

	comment, err := ctr.store.CreateComment(c.Context(), &model.Comment{
		PostID:   postID,
		Content:  p.Content,
		PostedBy: user.ID,
	})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(comment)
}
