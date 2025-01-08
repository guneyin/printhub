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
	AppName string `env:"PH_APP_NAME"`
	AppURL  string `env:"PH_APP_URL"`

	ApiBaseUrl string `env:"PH_API_BASE_URL"`
	ApiPort    string `env:"PH_API_PORT"`

	MailEnabled bool `env:"PH_MAIL_ENABLED"`

	JWTSecret string `env:"PH_JWT_SECRET"`
	JWTExp    string `env:"PH_JWT_EXP"`

	AuthSecret string `env:"PH_AUTH_SECRET"`

	EmailUser       string `env:"PH_EMAIL_USER"`
	EmailPassword   string `env:"PH_EMAIL_PASSWORD"`
	EmailSmtpServer string `env:"PH_EMAIL_SMTP_SERVER"`
	EmailSmtpPort   string `env:"PH_EMAIL_SMTP_PORT"`
	EmailSender     string `env:"PH_EMAIL_SENDER"`
	EmailIdentity   string `env:"PH_EMAIL_IDENTITY"`

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
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("error loading .env file")
	}

	cfg := &Config{
		AppName: os.Getenv("PH_APP_NAME"),
		AppURL:  os.Getenv("PH_APP_URL"),

		ApiBaseUrl: os.Getenv("PH_API_BASE_URL"),
		ApiPort:    os.Getenv("PH_API_PORT"),

		MailEnabled: os.Getenv("PH_MAIL_ENABLED") == "true",

		JWTSecret: os.Getenv("PH_JWT_SECRET"),
		JWTExp:    os.Getenv("PH_JWT_EXP"),

		AuthSecret: os.Getenv("PH_AUTH_SECRET"),

		EmailUser:       os.Getenv("PH_EMAIL_USER"),
		EmailPassword:   os.Getenv("PH_EMAIL_PASSWORD"),
		EmailSmtpServer: os.Getenv("PH_EMAIL_SMTP_SERVER"),
		EmailSmtpPort:   os.Getenv("PH_EMAIL_SMTP_PORT"),
		EmailSender:     os.Getenv("PH_EMAIL_SENDER"),
		EmailIdentity:   os.Getenv("PH_EMAIL_IDENTITY"),

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
