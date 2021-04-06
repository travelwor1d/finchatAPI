package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/upload"
	"github.com/finchatapp/finchat-api/internal/verify"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

var a *fiber.App

var userAuthToken, goatAuthToken string

func TestMain(m *testing.M) {
	appconfig.Init("../configs/config.yaml")

	verifySvc := verify.Mock{}

	db, err := store.Connect(appconfig.Config.MySQL)
	if err != nil {
		log.Printf("failed to connect to mySQL db: %v", err)
	}

	u := upload.Mock{}
	msg := messaging.Mock{}

	s := store.New(db)
	if err = seedDB(s); err != nil {
		log.Fatal(err)
	}
	jwtM := token.NewJWTManager(appconfig.Config.Auth.Secret, time.Duration(appconfig.Config.Auth.Duration)*time.Minute)
	ctr := controller.New(s, jwtM, verifySvc, u, msg)

	a = fiber.New()
	app.Setup(a, ctr)

	userAuthToken, err = login("user@gmail.com", "admin123")
	if err != nil {
		log.Fatal(err)
	}
	goatAuthToken, err = login("goat@gmail.com", "admin123")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func seedDB(s *store.Store) error {
	user, err := s.CreateUser(context.Background(), &model.User{
		FirstName: "Example",
		LastName:  "User",
		Phone:     "+489603962412",
		Email:     "user@gmail.com",
		Type:      "USER",
	}, "admin123")
	if err != nil {
		return err
	}
	code, err := s.CreateGoatInviteCode(context.Background(), user.ID)
	if err != nil {
		return err
	}
	_, err = s.CreateUser(context.Background(), &model.User{
		FirstName: "Example",
		LastName:  "Goat",
		Phone:     "+489603962412",
		Email:     "goat@gmail.com",
		Type:      "GOAT",
	}, "admin123", code)
	if err != nil {
		return err
	}
	return nil
}

func login(email, password string) (string, error) {
	reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	req := httptest.NewRequest("POST", "/auth/v1/login", strings.NewReader(reqBody))
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.Test(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unsuccessful login: %s", body)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	return body.Token, err
}
