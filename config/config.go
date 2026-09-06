package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port       int
	Logpath    string
	Assentpath string
}

type DatabaseConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret     string
	Expirehour int
}

var CFG Config

func Load() error {
	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&CFG); err != nil {
		return err
	}

	if CFG.Server.Logpath == "" {
		CFG.Server.Logpath = "logs"
	}
	CFG.Server.Logpath = filepath.Join(CFG.Server.Logpath, fmt.Sprintf("%s.log", time.Now().Format("Monday_15")))

	return nil
}
