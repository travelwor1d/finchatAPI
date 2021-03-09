package appconfig

type Auth struct {
	Secret   string
	Duration int
}

type Twilio struct {
	SID   string
	Token string
}

type Stripe struct {
	Key string
}

type MySQL struct {
	ConnectionString string
}

type AppConfig struct {
	Port   string
	Auth   Auth   `mapstructure:"auth"`
	Twilio Twilio `mapstructure:"twilio"`
	Stripe Stripe `mapstructure:"stripe"`
	MySQL  MySQL  `mapstructure:"mysql"`
}
