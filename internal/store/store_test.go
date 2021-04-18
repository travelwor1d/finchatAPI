package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/gopher-lib/config"
)

func TestStore(t *testing.T) {
	var conf appconfig.AppConfig
	if err := config.LoadFile(&conf, "../../configs/config.yaml", "../../.env"); err != nil {
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

	t.Run("UpdateUser", func(t *testing.T) {
		user, err := s.UpdateUser(context.Background(), 10, nil, nil, nil)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("user: %#v\n", user)
	})
}

func TestErrors(t *testing.T) {
	if err := fmt.Errorf("function F: %w", ErrNotFound); !errors.Is(err, ErrNotFound) {
		t.Errorf("ErrNotFound was not found inside err")
	}
}
