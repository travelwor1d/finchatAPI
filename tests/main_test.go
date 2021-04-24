package tests

import (
	"log"
	"os"
	"testing"

	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/token"
	"github.com/finchatapp/finchat-api/internal/upload"
	"github.com/finchatapp/finchat-api/internal/verify"
	"github.com/gofiber/fiber/v2"
)

var a *fiber.App

func TestMain(m *testing.M) {
	appconfig.Init("../configs/config.yaml")

	verifySvc := verify.Mock{}

	db, err := store.Connect(appconfig.Config.MySQL)
	if err != nil {
		log.Printf("failed to connect to mySQL db: %v", err)
	}

	u := upload.Mock{}
	msg := messaging.Mock{}
	tokenSvc := token.Mock{}

	s := store.New(db)
	if err = seedDB(s); err != nil {
		log.Fatal(err)
	}
	ctr := controller.New(s, tokenSvc, verifySvc, u, msg)

	a = fiber.New()
	app.Setup(a, ctr)

	os.Exit(m.Run())
}

func seedDB(s *store.Store) error {
	return nil
}
