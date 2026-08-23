package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

func Load() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
