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
	apiv1.Post("/goats/subscription-plans", ctr.CreateSubscriptionPlan)
	apiv1.Post("/users/credit-card", ctr.AddCreditCard)
	apiv1.Post("/users/subscriptions", ctr.CreateSubscription)

	apiv1.Get("/users", ctr.ListUsers)
	apiv1.Post("/users/me/profile-avatar:upload", ctr.UploadProvileAvatar)
	apiv1.Get("/users/:id", ctr.GetUser)
	apiv1.Patch("/users/:id", ctr.UpdateUser)
	apiv1.Delete("/users/:id", ctr.SoftDeleteUser)
	apiv1.Post("/users/:id:undelete", ctr.UndeleteUser)

	apiv1.Get("/posts", ctr.ListPosts)
	apiv1.Get("/posts/:id", ctr.GetPost)
	apiv1.Post("/posts", ctr.CreatePost)
	apiv1.Get("/comments", ctr.ListComments)
	apiv1.Get("/comments/:id", ctr.GetComment)
	apiv1.Post("/comments", ctr.CreateComment)

	// Post upvoting/downvoting
	apiv1.Post("/posts/:id/upvote", ctr.UpvotePost)
	apiv1.Delete("/posts/:id/upvote", ctr.RevertPostUpvote)
	apiv1.Post("/posts/:id/downvote", ctr.DownvotePost)
	apiv1.Delete("/posts/:id/downvote", ctr.RevertPostDownvote)
}
