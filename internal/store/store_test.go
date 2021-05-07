package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/model"
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
		user, err := s.UpdateUser(context.Background(), 10, nil, nil, nil, nil)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("user: %#v\n", user)
	})

	t.Run("DeleteUserByEmail", func(t *testing.T) {
		err := s.DeleteUserByEmail(context.Background(), "martilukas7@gmail.com")
		if err != nil {
			t.Error(err)
		}
		fmt.Println("deleted")
	})

	t.Run("CreateUser", func(t *testing.T) {
		_, err := s.CreateUser(context.Background(), &model.User{
			FirstName:   "Example",
			LastName:    "User",
			Phonenumber: "+489603962412",
			CountryCode: "PL",
			Email:       "user@gmail.com",
			Type:        "USER",
		}, "admin123")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestErrors(t *testing.T) {
	if err := fmt.Errorf("function F: %w", ErrNotFound); !errors.Is(err, ErrNotFound) {
		t.Errorf("ErrNotFound was not found inside err")
	}
}
