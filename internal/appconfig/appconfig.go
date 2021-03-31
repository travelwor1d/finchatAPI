package appconfig

import (
	"log"

	"github.com/gopher-lib/config"
)

var Config AppConfig

func init() {
	if err := config.LoadFile(&Config, "configs/config.yaml"); err != nil {
		log.Fatalf("failed to load app configuration: %v", err)
	}
}

type Auth struct {
	Secret   string
	Duration int
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
	Port    string
	Logger  bool
	Auth    Auth    `mapstructure:"auth"`
	Twilio  Twilio  `mapstructure:"twilio"`
	Stripe  Stripe  `mapstructure:"stripe"`
	Pubnub  Pubnub  `mapstructure:"pubnub"`
	MySQL   MySQL   `mapstructure:"mysql"`
	Storage Storage `mapstructure:"storage"`
}
