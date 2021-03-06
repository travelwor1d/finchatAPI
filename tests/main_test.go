package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/config"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

var a *fiber.App

func TestMain(m *testing.M) {
	var conf appconfig.AppConfig
	if err := config.LoadConfig(&conf, "../configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	db, err := store.Connect(conf.MySQL)
	if err != nil {
		log.Printf("failed to connect to mySQL db: %v", err)
	}

	s := store.New(db)
	jwtM := token.NewJWTManager(conf.Auth.Secret, time.Duration(conf.Auth.Duration)*time.Minute)
	ctr := controller.New(s, jwtM)

	a = fiber.New()
	app.Setup(a, ctr)

	os.Exit(m.Run())
}
