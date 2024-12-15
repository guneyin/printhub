package disk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/service/disk"
	"sync"
)

const handlerName = "disk"

type Handler struct {
	svc *disk.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: disk.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	//g := h.router.Group(h.name())
}
