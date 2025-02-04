package api

import (
	"encoding/gob"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/server"
	"github.com/guneyin/printhub/utils"
)

const appName = "PrintHub"

type Application struct {
	Name    string
	Version string
	Server  *fiber.App
	Handler *handler.Handler
}

func NewApplication() (*Application, error) {
	gob.Register(&model.Session{})

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
