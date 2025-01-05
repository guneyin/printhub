package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/admin"
	"github.com/guneyin/printhub/handler/auth"
	"github.com/guneyin/printhub/handler/config"
	"github.com/guneyin/printhub/handler/disk"
	"github.com/guneyin/printhub/handler/tenant"
	"github.com/guneyin/printhub/handler/user"
)

type IHandler interface {
	name() string
	setRoutes()
}

type Handler struct {
	router fiber.Router
}

func New(app *fiber.App) *Handler {
	handler := &Handler{
		router: app.Group("/api"),
	}
	handler.registerHandlers()

	return handler
}

func (h *Handler) registerHandlers() {
	admin.InitHandler(h.router)
	auth.InitHandler(h.router)
	user.InitHandler(h.router)
	config.InitHandler(h.router)
	tenant.InitHandler(h.router)
	disk.InitHandler(h.router)
}
