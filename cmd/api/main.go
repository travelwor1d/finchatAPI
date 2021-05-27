package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/logerr"
	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/token"
	"github.com/finchatapp/finchat-api/internal/upload"
	"github.com/finchatapp/finchat-api/internal/verify"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kevinburke/twilio-go"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"

	contactrepo "github.com/finchatapp/finchat-api/internal/entities/contact/repositories"
	contactusecase "github.com/finchatapp/finchat-api/internal/entities/contact/usecase"

	userrepo "github.com/finchatapp/finchat-api/internal/entities/user/repositories"
	userusecase "github.com/finchatapp/finchat-api/internal/entities/user/usecase"
)

func main() {
	// Initialize configuration service.
	appconfig.Init("configs/config.yaml")

	// Configure Twilio client.
	c := twilio.NewClient(appconfig.Config.Twilio.SID, appconfig.Config.Twilio.Token, nil).Verify.Verifications
	verifySvc := verify.New(c, appconfig.Config.Twilio.Verify)

	// Configure Stripe client.
	stripe.Key = appconfig.Config.Stripe.Key

	// Configure Firebase Authentication client.
	firebaseapp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed initialize firebase app: %v", err)
	}
	auth, err := firebaseapp.Auth(context.Background())
	if err != nil {
		log.Fatalf("failed initialize firebase auth client: %v", err)
	}
	tokenSvc := token.NewService(auth)

	// Configure Cloud Storage and upload service.
	storageClint, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create storage service: %v", err)
	}
	bkt := storageClint.Bucket(appconfig.Config.Storage.BucketName)
	u := upload.New(bkt)

	// Connect to mySQL database and configure store service.
	db, err := store.Connect(appconfig.Config.MySQL)
	if err != nil {
		log.Fatalf("failed to connect to mySQL db: %v", err)
	}
	s := store.New(db)

	// Configure messaging service.
	msg := messaging.New(appconfig.Config.Pubnub, s)

	projectID := "finchat-api-staging"
	errorClient, err := errorreporting.NewClient(context.Background(), projectID, errorreporting.Config{
		ServiceName: "finchat-api",
		OnError: func(err error) {
			log.Fatal(err)
		},
	})
	if err != nil {
		log.Fatalf("failed to create error reporting client: %v", err)
	}
	defer errorClient.Close()

	lr := logerr.New(errorClient)

	// tmp code injection
	dsn := appconfig.Config.MySQL.ConnectionString
	master, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logrus.Fatal("master db connection", err)
	}
	defer master.Close()

	dbs := []*sqlx.DB{master}
	contactuc := contactusecase.New(contactrepo.New(dbs))
	useruc := userusecase.New(userrepo.New(dbs))
	// tmp code injection

	// Setup application controller with its dependencies.
	ctr := controller.New(s, contactuc, useruc, tokenSvc, verifySvc, u, msg, lr)

	// Configure and run fiber.
	a := fiber.New()
	app.Setup(a, ctr)
	addr := fmt.Sprintf(":%s", appconfig.Config.Port)
	log.Fatal(a.Listen(addr))
}
