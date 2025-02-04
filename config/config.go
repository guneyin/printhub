package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string `env:"PH_APP_NAME"`
	AppURL  string `env:"PH_APP_URL"`

	APIBaseURL string `env:"PH_API_BASE_URL"`
	APIPort    string `env:"PH_API_PORT"`

	MailEnabled bool `env:"PH_MAIL_ENABLED"`

	JWTSecret string `env:"PH_JWT_SECRET"`
	JWTExp    string `env:"PH_JWT_EXP"`

	AuthSecret string `env:"PH_AUTH_SECRET"`

	EmailUser       string `env:"PH_EMAIL_USER"`
	EmailPassword   string `env:"PH_EMAIL_PASSWORD"`
	EmailSMTPServer string `env:"PH_EMAIL_SMTP_SERVER"`
	EmailSMTPPort   string `env:"PH_EMAIL_SMTP_PORT"`
	EmailSender     string `env:"PH_EMAIL_SENDER"`
	EmailIdentity   string `env:"PH_EMAIL_IDENTITY"`

	DBPath string `env:"PH_DB_PATH"`
	DBHost string `env:"PH_DB_HOST"`
	DBPort string `env:"PH_DB_PORT"`
	DBUser string `env:"PH_DB_USER"`
	DBPwd  string `env:"PH_DB_PWD"`
	DBName string `env:"PH_DB_NAME"`

	GoogleClientID     string `env:"PH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"PH_GOOGLE_CLIENT_SECRET"`
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("error loading .env file")
	}

	cfg := &Config{
		AppName: os.Getenv("PH_APP_NAME"),
		AppURL:  os.Getenv("PH_APP_URL"),

		APIBaseURL: os.Getenv("PH_API_BASE_URL"),
		APIPort:    os.Getenv("PH_API_PORT"),

		MailEnabled: os.Getenv("PH_MAIL_ENABLED") == "true",

		JWTSecret: os.Getenv("PH_JWT_SECRET"),
		JWTExp:    os.Getenv("PH_JWT_EXP"),

		AuthSecret: os.Getenv("PH_AUTH_SECRET"),

		EmailUser:       os.Getenv("PH_EMAIL_USER"),
		EmailPassword:   os.Getenv("PH_EMAIL_PASSWORD"),
		EmailSMTPServer: os.Getenv("PH_EMAIL_SMTP_SERVER"),
		EmailSMTPPort:   os.Getenv("PH_EMAIL_SMTP_PORT"),
		EmailSender:     os.Getenv("PH_EMAIL_SENDER"),
		EmailIdentity:   os.Getenv("PH_EMAIL_IDENTITY"),

		DBPath: os.Getenv("PH_DB_PATH"),
		DBHost: "",
		DBPort: "",
		DBUser: "",
		DBPwd:  "",
		DBName: "",

		GoogleClientID:     os.Getenv("PH_GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("PH_GOOGLE_CLIENT_SECRET"),
	}
	return cfg.validate()
}

func (c *Config) validate() (*Config, error) {
	if c.APIPort == "" {
		c.APIPort = "8080"
	}
	if c.APIBaseURL == "" {
		c.APIBaseURL = fmt.Sprintf("http://localhost:%s", c.APIPort)
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
