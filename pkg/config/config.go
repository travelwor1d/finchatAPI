package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig(rawVal interface{}, filename string, envPath ...string) error {
	if len(envPath) > 0 && envPath[0] != "" {
		if err := godotenv.Load(envPath[0]); err != nil {
			return fmt.Errorf("cannot load env. file: %v", err)
		}
	} else {
		// Ignore error as we are loading default env. file.
		_ = godotenv.Load(".env")
	}

	viper.SetConfigName(strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename))))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Dir(filename))
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return fmt.Errorf("cannot read in config: %w", err)
	}

	replaceEnv(viper.AllKeys())

	err = viper.Unmarshal(&rawVal)
	if err != nil {
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}
	return nil
}

var validEnv = regexp.MustCompile(`^\$\{[a-zA-Z_]+[a-zA-Z0-9_]*\}$`)

func replaceEnv(keys []string) {
	for _, key := range keys {
		old := viper.GetString(key)
		var new string
		if validEnv.MatchString(old) {
			new = os.Getenv(old[2 : len(old)-1])
		} else {
			new = old
		}
		viper.Set(key, new)
	}
}
