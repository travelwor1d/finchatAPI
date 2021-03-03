package store

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/pkg/config"
)

func TestStore(t *testing.T) {
	var conf appconfig.AppConfig
	if err := config.LoadConfig(&conf, "../../configs/config.yaml", "../../.env"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	db, err := Connect(conf.MySQL)
	if err != nil {
		log.Printf("failed to connect to mySQL db: %v", err)
	}
	s := New(db)

	t.Run("GetUserByEmail", func(t *testing.T) {
		user, err := s.GetUserByEmail(context.Background(), "martilukas7@gmail.com")
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("user: %#v\n", user)
	})
}
