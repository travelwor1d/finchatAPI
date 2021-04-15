package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

func MustParseClaims(jm *token.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idToken := c.Cookies("token")
		if idToken == "" {
			auth := c.Get("Authorization")
			if auth == "" {
				return httperr.New(codes.Omit, http.StatusUnauthorized, "missing bearer token").Send(c)
			}
			splitAuth := strings.Split(auth, "Bearer ")
			if len(splitAuth) != 2 {
				return httperr.New(codes.Omit, http.StatusUnauthorized, "malformed bearer token").Send(c)
			}
			idToken = splitAuth[1]
		}
		claims, err := jm.Verify(idToken)
		if errors.Is(err, token.ErrExpired) && claims != nil {
			return httperr.New(codes.Omit, http.StatusUnauthorized, "expired JWT").
				SetDetail(fmt.Sprintf("JWT expired at: %v", time.Unix(claims.ExpiresAt, 0))).Send(c)
		}
		if err != nil {
			return httperr.New(codes.Omit, http.StatusUnauthorized, "invalid JWT").Send(c)
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}
