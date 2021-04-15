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
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

type LimiterConfig struct {
	Max      int
	Duration time.Duration
}

func Limiter(config *LimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:      config.Max,
		Duration: config.Duration,
		Key: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return httperr.New(codes.Omit, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)).Send(c)
		},
	})
}
