package appconfig

type Auth struct {
	Secret   string
	Duration int
}

type MySQL struct {
	ConnectionString string
}

type AppConfig struct {
	Port  string
	Auth  `mapstructure:"auth"`
	MySQL `mapstructure:"mysql"`
}
