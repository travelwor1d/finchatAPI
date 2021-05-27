package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

var (
	errInternal  = httperr.New(codes.Omit, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	errParseBody = httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body")
)

func (ctr *Ctr) userFromCtx(c *fiber.Ctx) (*model.User, *httperr.HTTPErr) {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return nil, httperr.New(codes.Omit, http.StatusUnauthorized, "failed to retrieve uid from request context")
	}
	user, err := ctr.store.GetUserByFirebaseID(c.Context(), uid)
	if err != nil {
		return nil, httperr.New(codes.Omit, http.StatusUnauthorized, "User was not found in Firebase. Please try again").SetDetail(err)
	}
	if user != nil {
		user.LastSeen = time.Now().UTC().Truncate(time.Second)
		err = ctr.store.UpdateLastSeenUser(c.Context(), user.ID, user.LastSeen)
		if err != nil {
			return nil, httperr.New(codes.Omit, http.StatusBadRequest, "User was not updated. Please try again").SetDetail(err)
		}
	}
	return user, nil
}

func getUserTypes(t string) []string {
	if t == "" {
		return []string{"GOAT", "USER"}
	}
	types := strings.Split(t, ",")
	res := make([]string, 0)
	for _, typ := range types {
		typ = strings.ToUpper(strings.TrimSpace(typ))
		res = append(res, typ)
	}
	return res
}

func sendSuccess(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}
