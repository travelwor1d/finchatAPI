package appconfig

type MySQL struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

type AppConfig struct {
	Port  string
	MySQL `mapstructure:"mysql"`
}
