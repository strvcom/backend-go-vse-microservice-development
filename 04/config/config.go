package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port int
}

func ReadConfigFromFile(path string) (Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("reading configuration: %w", err)
	}

	cfg := Config{
		Port: viper.GetInt("port"),
	}

	return cfg, nil
}
