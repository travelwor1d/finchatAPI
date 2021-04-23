package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/finchatapp/finchat-api/internal/token"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func MustParseClaims(tokenSvc token.Service) fiber.Handler {
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

		token, err := tokenSvc.VerifyIDToken(c.Context(), idToken)
		if err != nil {
			return httperr.New(codes.Omit, http.StatusUnauthorized, "invalid JWT").SetDetail(err).Send(c)
		}

		c.Locals("uid", token.UID)
		return c.Next()
	}
}

type LimiterConfig struct {
	Max        int
	Expiration time.Duration
}

func Limiter(config *LimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return httperr.New(codes.Omit, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)).Send(c)
		},
	})
}
