package app

import (
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App, ctr *controller.Ctr) {
	p := middleware.Protected(ctr.JWTManager())
	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	authv1 := app.Group("/auth/v1")
	authv1.Post("/login", ctr.Login)
	authv1.Post("/register", ctr.Register)
	authv1.Get("/verify", p, ctr.RequestVerification)
	authv1.Post("/verify", p, ctr.Verify)

	apiv1 := app.Group("/api/v1", p)
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("v1")
	})
	apiv1.Get("/goats:invite", ctr.InviteGoat)
	apiv1.Get("/goats/invite-codes/:inviteCode", ctr.VerifyInviteCode)
}
