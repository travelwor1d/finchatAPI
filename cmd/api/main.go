package main

import (
	"fmt"
	"log"
	"time"

	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gopher-lib/config"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
)

func main() {
	var conf appconfig.AppConfig
	if err := config.LoadFile(&conf, "configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	verify := twilio.NewClient(conf.Twilio.SID, conf.Twilio.Token, nil).Verify.Verifications

	stripe.Key = conf.Stripe.Key

	db, err := store.Connect(conf.MySQL)
	if err != nil {
		log.Fatalf("failed to connect to mySQL db: %v", err)
	}

	s := store.New(db)
	jwtM := token.NewJWTManager(conf.Auth.Secret, time.Duration(conf.Auth.Duration)*time.Minute)
	ctr := controller.New(s, jwtM, verify)

	a := fiber.New()

	app.Setup(a, ctr)

	addr := fmt.Sprintf(":%s", conf.Port)
	log.Fatal(a.Listen(addr))
}
