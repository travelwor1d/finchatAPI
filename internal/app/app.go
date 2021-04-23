package app

import (
	"time"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App, ctr *controller.Ctr) {
	p := middleware.MustParseClaims(ctr.TokenSvc())
	l := middleware.Limiter(&middleware.LimiterConfig{Max: 5, Expiration: time.Second})
	// Global middleware
	app.Use(recover.New())
	if appconfig.Config.Logger {
		app.Use(logger.New())
	}
	app.Use(cors.New())

	authv1 := app.Group("/auth/v1")
	authv1.Post("/register", ctr.Register)
	authv1.Get("/verify", p, ctr.RequestVerification)
	authv1.Post("/verify", p, ctr.Verify)
	authv1.Get("/email-validation", l, ctr.EmailValidation)
	authv1.Get("/phone-number-validation", l, ctr.PhonenumberValidation)

	apiv1 := app.Group("/api/v1", p)
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("v1")
	})
	// TODO: use \\
	apiv1.Get("/users/goats:invite", ctr.InviteGoat)
	apiv1.Get("/users/goats/invite-codes/:inviteCode", ctr.VerifyInviteCode)
	apiv1.Post("/users/goats/subscription-plans", ctr.CreateSubscriptionPlan)
	apiv1.Post("/users/credit-card", ctr.AddCreditCard)
	apiv1.Post("/users/subscriptions", ctr.CreateSubscription)

	apiv1.Get("/users", ctr.ListUsers)
	apiv1.Post("/users/me/profile-avatar", ctr.UploadProvileAvatar)
	apiv1.Get("/users/:id", ctr.GetUser)
	apiv1.Get("/users/me", ctr.GetUser)
	apiv1.Patch("/users/me", ctr.UpdateUser)
	apiv1.Delete("/users/me", ctr.SoftDeleteUser)

	apiv1.Get("/users/me/contacts", ctr.ListContactsOrRequests)
	apiv1.Get("/users/me/contacts/:id", ctr.GetContact)
	apiv1.Post("/users/me/contacts", ctr.CreateContactRequest)
	apiv1.Patch("/users/me/contacts/:id", ctr.PatchContactRequest)
	apiv1.Delete("/users/me/contacts/:id", ctr.DeleteContact)

	// Blogging API
	apiv1.Get("/posts", ctr.ListPosts)
	apiv1.Get("/posts/:id", ctr.GetPost)
	apiv1.Post("/posts", ctr.CreatePost)
	apiv1.Get("/posts/:postId/comments", ctr.ListComments)
	apiv1.Get("/posts/:postId/comments/:id", ctr.GetComment)
	apiv1.Post("/posts/:postId/comments", ctr.CreateComment)

	// Post upvoting/downvoting
	apiv1.Post("/posts/:id/upvote", ctr.UpvotePost)
	apiv1.Delete("/posts/:id/upvote", ctr.RevertPostUpvote)
	apiv1.Post("/posts/:id/downvote", ctr.DownvotePost)
	apiv1.Delete("/posts/:id/downvote", ctr.RevertPostDownvote)
}
