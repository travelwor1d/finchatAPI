package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("simple config with interpolation", func(t *testing.T) {
		type Config struct {
			Port    int
			Secret1 string
			Secret2 string
		}
		var conf Config
		err := LoadConfig(&conf, "testdata/config.testing.yaml", "testdata/.env.testing")
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{5432, "secret-value", ""}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})

	t.Run("more complex config with interpolation", func(t *testing.T) {
		type DB struct {
			User     string
			Password string
		}
		type Config struct {
			Port int
			DB   DB
		}
		var conf Config
		err := LoadConfig(&conf, "testdata/config.testing.yaml", "testdata/.env.testing")
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{5432, DB{"admin", "root"}}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})
}
