package server

import (
	"errors"
	"fmt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/guneyin/printhub/handler/mw"
	"os"
	"time"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
)

func NewServer(appName string) *fiber.App {
	app := fiber.New(fiber.Config{
		ServerHeader:      fmt.Sprintf("%s HTTP Server", appName),
		BodyLimit:         16 * 1024 * 1024,
		AppName:           appName,
		EnablePrintRoutes: true,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return ctx.Status(code).JSON(mw.HTTPError{
				Error: err.Error(),
			})
		},
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(favicon.New())

	if _, err := os.Stat("./docs/swagger.json"); err == nil {
		app.Use(swagger.New(swagger.Config{
			BasePath: "/api/",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API Docs",
		}))
	}

	app.Get("/test", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "ok"})
	})

	return app
}
