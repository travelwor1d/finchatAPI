package controller

import (
	"net/http"
	"strconv"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

func userID(c *fiber.Ctx) (int, *httperr.HTTPErr) {
	claims, ok := c.Locals("claims").(*token.JWTClaims)
	if !ok || claims == nil {
		return 0, httperr.New(codes.Omit, http.StatusUnauthorized, "failed to retrieve claims from request context")
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, httperr.New(codes.Omit, http.StatusInternalServerError, "invalid claims subject")
	}
	return id, nil
}
