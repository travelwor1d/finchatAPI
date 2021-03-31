package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/upload"
	"github.com/finchatapp/finchat-api/internal/verify"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
)

func main() {
	c := twilio.NewClient(appconfig.Config.Twilio.SID, appconfig.Config.Twilio.Token, nil).Verify.Verifications
	verifySvc := verify.New(c, appconfig.Config.Twilio.Verify)

	stripe.Key = appconfig.Config.Stripe.Key

	db, err := store.Connect(appconfig.Config.MySQL)
	if err != nil {
		log.Fatalf("failed to connect to mySQL db: %v", err)
	}

	storageClint, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create storage service: %v", err)
	}
	bkt := storageClint.Bucket(appconfig.Config.Storage.BucketName)
	u := upload.New(bkt)

	s := store.New(db)
	jwtM := token.NewJWTManager(appconfig.Config.Auth.Secret, time.Duration(appconfig.Config.Auth.Duration)*time.Minute)
	msg := messaging.New(appconfig.Config.Pubnub, s)
	ctr := controller.New(s, jwtM, verifySvc, u, msg)

	a := fiber.New()

	app.Setup(a, ctr)

	addr := fmt.Sprintf(":%s", appconfig.Config.Port)
	log.Fatal(a.Listen(addr))
}
