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
	TLS      TLSConfig
}

type ServerConfig struct {
	Port       int
	Logpath    string
	Assentpath string
	Frontend   FrontendConfig
}

type FrontendConfig struct {
	Enabled bool
	Path    string
}

type TLSConfig struct {
	Enabled bool
	Cert    string
	Key     string
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

	if CFG.Server.Frontend.Path == "" {
		CFG.Server.Frontend.Path = "frontend/dist"
	}

	if CFG.TLS.Cert == "" {
		CFG.TLS.Cert = "cert.pem"
	}
	if CFG.TLS.Key == "" {
		CFG.TLS.Key = "key.pem"
	}

	return nil
}
