package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	ApiBaseUrl string `env:"PH_API_BASE_URL"`
	ApiPort    string `env:"PH_API_PORT"`

	DbPath string `env:"PH_DB_PATH"`
	DbHost string `env:"PH_DB_HOST"`
	DbPort string `env:"PH_DB_PORT"`
	DbUser string `env:"PH_DB_USER"`
	DbPwd  string `env:"PH_DB_PWD"`
	DbName string `env:"PH_DB_NAME"`

	GoogleClientId     string `env:"PH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"PH_GOOGLE_CLIENT_SECRET"`
}

func New() *Config {
	err := godotenv.Load("test.env")
	if err != nil {
		slog.Warn("error loading .env file")
	}

	cfg := &Config{
		ApiBaseUrl:         os.Getenv("PH_API_BASE_URL"),
		ApiPort:            os.Getenv("PH_API_PORT"),
		DbPath:             os.Getenv("PH_DB_PATH"),
		DbHost:             "",
		DbPort:             "",
		DbUser:             "",
		DbPwd:              "",
		DbName:             "",
		GoogleClientId:     os.Getenv("PH_GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("PH_GOOGLE_CLIENT_SECRET"),
	}
	_ = cfg.validate()
	return cfg
}

func (c *Config) validate() error {
	if c.ApiPort == "" {
		c.ApiPort = "8080"
	}
	if c.ApiBaseUrl == "" {
		c.ApiBaseUrl = fmt.Sprintf("http://localhost:%s", c.ApiPort)
	}

	return nil
}
