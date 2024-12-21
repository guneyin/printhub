package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ApiBaseUrl string `env:"PH_API_BASE_URL"`
	ApiPort    string `env:"PH_API_PORT"`

	JWTSecret string `env:"PH_JWT_SECRET"`
	JWTExp    string `env:"PH_JWT_EXP"`

	DbPath string `env:"PH_DB_PATH"`
	DbHost string `env:"PH_DB_HOST"`
	DbPort string `env:"PH_DB_PORT"`
	DbUser string `env:"PH_DB_USER"`
	DbPwd  string `env:"PH_DB_PWD"`
	DbName string `env:"PH_DB_NAME"`

	GoogleClientId     string `env:"PH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"PH_GOOGLE_CLIENT_SECRET"`
}

func New() (*Config, error) {
	err := godotenv.Load("test.env")
	if err != nil {
		slog.Warn("error loading .env file")
	}

	cfg := &Config{
		ApiBaseUrl: os.Getenv("PH_API_BASE_URL"),
		ApiPort:    os.Getenv("PH_API_PORT"),

		JWTSecret: os.Getenv("PH_JWT_SECRET"),
		JWTExp:    os.Getenv("PH_JWT_EXP"),

		DbPath: os.Getenv("PH_DB_PATH"),
		DbHost: "",
		DbPort: "",
		DbUser: "",
		DbPwd:  "",
		DbName: "",

		GoogleClientId:     os.Getenv("PH_GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("PH_GOOGLE_CLIENT_SECRET"),
	}
	return cfg.validate()
}

func (c *Config) validate() (*Config, error) {
	if c.ApiPort == "" {
		c.ApiPort = "8080"
	}
	if c.ApiBaseUrl == "" {
		c.ApiBaseUrl = fmt.Sprintf("http://localhost:%s", c.ApiPort)
	}
	if strings.TrimSpace(c.JWTSecret) == "" {
		return nil, errors.New("JWT secret is required")
	}

	exp, _ := strconv.Atoi(c.JWTExp)
	if exp == 0 {
		exp = 30
	}
	c.JWTExp = strconv.Itoa(exp)

	return c, nil
}
