package tests

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/model"
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
	users := []*model.User{
		{FirebaseID: "123456", FirstName: "John", LastName: "Doe", Email: "john.doe@gmail.com", Phonenumber: "48907689431", CountryCode: "PL"},
		{FirebaseID: "654321", FirstName: "Jane", LastName: "Doe", Email: "jane.doe@gmail.com", Phonenumber: "48907689432", CountryCode: "PL"},
	}

	var err error
	wg := sync.WaitGroup{}
	wg.Add(len(users))
	for _, user := range users {
		go func(u *model.User) {
			_, err = s.CreateUser(context.Background(), u)
			wg.Done()
		}(user)
	}
	wg.Wait()

	return err
}
