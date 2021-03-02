package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/config"
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
	_ = s

	app := fiber.New()

	addr := fmt.Sprintf(":%s", conf.Port)
	log.Fatal(app.Listen(addr))
}
