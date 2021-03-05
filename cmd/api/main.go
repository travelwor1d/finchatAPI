package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/middleware"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/config"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	var conf appconfig.AppConfig
	if err := config.LoadConfig(&conf, "configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	firebaseapp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to initialize firebase app: %v", err)
	}

	auth, err := firebaseapp.Auth(context.Background())
	if err != nil {
		log.Fatalf("failed to initialize firebase auth: %v", err)
	}
	_ = auth

	db, err := store.Connect(conf.MySQL)
	if err != nil {
		log.Printf("failed to connect to mySQL db: %v", err)
	}
	s := store.New(db)

	jwtM := token.NewJWTManager(conf.Auth.Secret, time.Duration(conf.Auth.Duration)*time.Minute)

	ctr := controller.New(s, jwtM)

	app := fiber.New()
	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	authv1 := app.Group("/auth/v1")
	authv1.Post("/login", ctr.Login)
	authv1.Post("/register", ctr.Register)

	apiv1 := app.Group("/api/v1", middleware.Protected(jwtM))
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("v1")
	})
	apiv1.Get("/goats:invite", ctr.InviteGoat)
	apiv1.Get("/goats/invite-codes/:inviteCode", ctr.VerifyInviteCode)

	addr := fmt.Sprintf(":%s", conf.Port)
	log.Fatal(app.Listen(addr))
}
