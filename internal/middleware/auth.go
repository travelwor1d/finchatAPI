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
		idToken := c.Cookies("token")
		if idToken == "" {
			auth := c.Get("Authorization")
			splitAuth := strings.Split(auth, "Bearer ")
			if len(splitAuth) != 2 {
				return httperr.New(codes.Omit, fiber.StatusBadRequest, "missing or malformed JWT").Send(c)
			}
			idToken = splitAuth[1]
		}
		claims, err := jm.Verify(idToken)
		if err != nil {
			return httperr.New(codes.Omit, fiber.StatusUnauthorized, "invalid or expired JWT").Send(c)
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}
