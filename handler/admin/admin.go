package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/admin"
	"sync"
)

const handlerName = "admin"

type Handler struct {
	svc *admin.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: admin.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.Group(h.name())

	g.Post("/tenant", h.tenantCreate)
	tenant := g.Group("/tenant/:id")

	user := tenant.Group("/user")
	user.Post("/", h.tenantUserCreate)
}

func (h *Handler) tenantCreate(c *fiber.Ctx) error {
	tenant, err := model.NewTenant(c.Body())
	if err != nil {
		return err
	}
	err = h.svc.TenantCreate(c.Context(), tenant)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(tenant)
}

func (h *Handler) tenantUserCreate(c *fiber.Ctx) error {
	return nil
}
