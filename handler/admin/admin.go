package admin

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/admin"
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
	g := r.Group(h.name()).Use(mw.AdminGuard)

	tenant := g.Group("/tenant")
	tenant.Post("/", h.saveTenant)
	tenant.Get("/", h.searchTenant)
	tenant.Get("/id", h.getTenantByID)
	tenant.Delete("/id", h.deleteTenant)
	tenant.Post("/user", h.addTenantUser)
	tenant.Get("/user", h.getTenantUserList)
}

// @Router /admin/tenant [post].
func (h *Handler) saveTenant(c *fiber.Ctx) error {
	tenant, err := model.NewTenant(c.Body())
	if err != nil {
		return err
	}
	err = h.svc.SaveTenant(c.Context(), tenant)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(tenant)
}

func (h *Handler) searchTenant(c *fiber.Ctx) error {
	filter := new(model.TenantFilter)
	if err := c.QueryParser(filter); err != nil {
		return err
	}

	tenant, err := h.svc.SearchTenant(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(tenant)
}

func (h *Handler) getTenantByID(c *fiber.Ctx) error {
	id := c.Params("id")

	tenant, err := h.svc.GetTenantByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(tenant)
}

func (h *Handler) deleteTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.DeleteTenant(c.Context(), id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) addTenantUser(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	userEmail := c.Params("user_email")

	u := &model.User{Role: model.UserRoleTenant, Email: userEmail}
	err := h.svc.AddTenantUser(c.Context(), tenantID, u)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) getTenantUserList(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	list, err := h.svc.GetTenantUserList(c.Context(), tenantID)
	if err != nil {
		return err
	}

	return c.JSON(list)
}
