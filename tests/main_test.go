package tests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/finchatapp/finchat-api/internal/app"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/controller"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gopher-lib/config"
	"github.com/kevinburke/twilio-go"
)

var a *fiber.App

var authToken string

func TestMain(m *testing.M) {
	var conf appconfig.AppConfig
	if err := config.LoadFile(&conf, "../configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}

	verify := twilio.NewClient(conf.Twilio.SID, conf.Twilio.Token, nil).Verify.Verifications

	db, err := store.Connect(conf.MySQL)
	if err != nil {
		log.Printf("failed to connect to mySQL db: %v", err)
	}

	s := store.New(db)
	jwtM := token.NewJWTManager(conf.Auth.Secret, time.Duration(conf.Auth.Duration)*time.Minute)
	ctr := controller.New(s, jwtM, verify)

	a = fiber.New()
	app.Setup(a, ctr)

	req := httptest.NewRequest("POST", "/auth/v1/login", strings.NewReader(`
	{
		"email": "example@gmail.com",
		"password": "admin123"
	}
	`))
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.Test(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("unsuccessful login: %s", body)
	}
	var body struct{ JWT string }
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}
	authToken = body.JWT

	os.Exit(m.Run())
}
