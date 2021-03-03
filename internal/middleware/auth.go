package middleware

import (
	"strings"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

func Protected(jm *token.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("token")
		if token == "" {
			auth := c.Get("Authorization")
			splitAuth := strings.Split(auth, "Bearer ")
			if len(splitAuth) != 2 {
				return httperr.New(codes.Omit, fiber.StatusBadRequest, "missing or malformed JWT").Send(c)
			}
			token = splitAuth[1]
		}
		_, err := jm.Verify(token)
		if err != nil {
			return httperr.New(codes.Omit, fiber.StatusUnauthorized, "invalid or expired JWT").Send(c)
		}
		return c.Next()
	}
}
