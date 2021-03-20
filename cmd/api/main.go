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
	"github.com/gopher-lib/config"
	"github.com/kevinburke/twilio-go"
	"github.com/stripe/stripe-go/v72"
)

func main() {
	var conf appconfig.AppConfig
	if err := config.LoadFile(&conf, "configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	c := twilio.NewClient(conf.Twilio.SID, conf.Twilio.Token, nil).Verify.Verifications
	verifySvc := verify.New(c, conf.Twilio.Verify)

	stripe.Key = conf.Stripe.Key

	db, err := store.Connect(conf.MySQL)
	if err != nil {
		log.Fatalf("failed to connect to mySQL db: %v", err)
	}

	msg := messaging.New(conf.Pubnub)

	storageClint, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	bkt := storageClint.Bucket(conf.Storage.BucketName)
	u := upload.New(bkt)

	s := store.New(db)
	jwtM := token.NewJWTManager(conf.Auth.Secret, time.Duration(conf.Auth.Duration)*time.Minute)
	ctr := controller.New(s, jwtM, verifySvc, u, msg)

	a := fiber.New()

	app.Setup(a, ctr)

	addr := fmt.Sprintf(":%s", conf.Port)
	log.Fatal(a.Listen(addr))
}
