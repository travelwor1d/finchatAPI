package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/config"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
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

	a := fiber.New()

	app.Setup(a, ctr)

	addr := fmt.Sprintf(":%s", conf.Port)
	log.Fatal(a.Listen(addr))
}
