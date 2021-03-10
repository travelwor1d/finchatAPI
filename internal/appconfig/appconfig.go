package appconfig

type Auth struct {
	Secret   string
	Duration int
}

type Twilio struct {
	SID   string
	Token string
}

type MySQL struct {
	ConnectionString string
}

type AppConfig struct {
	Port   string
	Auth   `mapstructure:"auth"`
	Twilio `mapstructure:"twilio"`
	MySQL  `mapstructure:"mysql"`
}
