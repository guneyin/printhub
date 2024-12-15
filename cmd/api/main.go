package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/server"
	"github.com/guneyin/printhub/utils"
	"log/slog"
	"os"
)

const appName = "PrintHub"

type Application struct {
	Name    string
	Version string
	Server  *fiber.App
	Handler *handler.Handler
}

func NewApplication() (*Application, error) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	market.InitMarket()

	appServer := server.NewServer(appName)
	appHandler := handler.New(appServer)

	return &Application{
		Name:    appName,
		Version: utils.GetVersion().Version,
		Server:  appServer,
		Handler: appHandler,
	}, nil
}
