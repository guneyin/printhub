package config

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/config"
)

const handlerName = "config"

type Handler struct {
	svc *config.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: config.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.Group(h.name())

	g.Get("/", h.get)
	g.Put("/", h.put)
	g.Delete("/", h.delete)
}

// @Router /config [get].
func (h *Handler) get(c *fiber.Ctx) error {
	// id := c.Locals("orgId").(string)
	id := c.Query("identity")
	cfg, err := h.svc.Get(c.Context(), id, c.Query("module"), c.Query("key"))
	if err != nil {
		return err
	}

	return c.JSON(cfg)
}

// @Router /config [put].
func (h *Handler) put(c *fiber.Ctx) error {
	cfg, err := model.NewConfigList(c.Body())
	if err != nil {
		return err
	}

	return h.svc.Set(c.Context(), cfg)
}

// @Router /config [delete].
func (h *Handler) delete(c *fiber.Ctx) error {
	id := c.Query("identity")
	return h.svc.Delete(c.Context(), id, c.Query("module"), c.Query("key"))
}
