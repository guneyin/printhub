package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/service/user"
	"sync"
)

const handlerName = "user"

type Handler struct {
	svc *user.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: user.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	//g := h.r.Group(h.name())
}
