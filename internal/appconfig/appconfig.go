package appconfig

import (
	"log"

	"github.com/gopher-lib/config"
)

var Config AppConfig

func Init(filename string) {
	if err := config.LoadFile(&Config, filename); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}
}

type Twilio struct {
	SID    string
	Token  string
	Verify string
}

type Stripe struct {
	Key string
}

type MySQL struct {
	ConnectionString string
}

type Storage struct {
	BucketName string
}

type Pubnub struct {
	PubKey     string
	SubKey     string
	SecKey     string
	ServerUUID string
}

type AppConfig struct {
	Port           string
	Logger         bool
	ErrorReporting bool
	WebhookToken   string
	Twilio         Twilio  `mapstructure:"twilio"`
	Stripe         Stripe  `mapstructure:"stripe"`
	Pubnub         Pubnub  `mapstructure:"pubnub"`
	MySQL          MySQL   `mapstructure:"mysql"`
	Storage        Storage `mapstructure:"storage"`
}
