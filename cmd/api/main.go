package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nkanders/finchat-api/internal/appconfig"
	"github.com/nkanders/finchat-api/pkg/config"
)

func main() {
	var conf appconfig.AppConfig
	if err := config.LoadConfig(&conf, "configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	app := fiber.New()

	addr := fmt.Sprintf(":%s", conf.Port)
	log.Fatal(app.Listen(addr))
}
