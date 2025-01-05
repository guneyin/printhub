package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/service/tenant"
	"sync"
)

const handlerName = "tenant"

type Handler struct {
	svc *tenant.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: tenant.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.Group(h.name()).Use(mw.AdminGuard)

	g.Get("/", h.getTenantList)
	//g.Get("/:id", h.getTenant)
}

func (h *Handler) getTenantList(c *fiber.Ctx) error {
	list, err := h.svc.GetTenantList(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(list)
}
