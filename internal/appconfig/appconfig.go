package appconfig

type MySQL struct {
	ConnectionString string
}

type AppConfig struct {
	Port  string
	MySQL `mapstructure:"mysql"`
}
