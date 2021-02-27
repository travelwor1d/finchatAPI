package main

import (
	"fmt"
	"log"

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
