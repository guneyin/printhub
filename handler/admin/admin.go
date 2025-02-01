package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
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
	g := r.Group(h.name()).Use(mw.AdminGuard)

	tenant := g.Group("/tenant")
	tenant.Get("/", h.getTenantList)
	tenant.Post("/", h.saveTenant)
	tenant.Get("/id", h.getTenantByID)
	tenant.Post("/user", h.saveTenantUser)
}

func (h *Handler) boostrap(c *fiber.Ctx) error {
	cnt, err := h.svc.GetCount(c.Context())
	if err != nil {
		return err
	}

	if cnt > 0 {
		return c.Next()
	}

	if c.Method() == fiber.MethodGet {
		return c.Redirect(
			"/admin/auth/register",
			fiber.StatusTemporaryRedirect)
	}

	return c.Next()
}

func (h *Handler) getTenantList(c *fiber.Ctx) error {
	//tenant, err := h.svc.GetTenantById(c.Context(), c.Query("filter"))
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

// saveTenant
// @Summary tenant create.
// @Description Create a new tenant.
// @Tags tenant create
// @Accept json
// @Produce json
// @Param tenant body model.Tenant true "tenant"
// @Failure default {object} mw.HTTPError
// @Router /admin/tenant [post]
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

func (h *Handler) saveTenantUser(c *fiber.Ctx) error {
	return nil
}

//func (h *Handler) getTenantUser(c *fiber.Ctx) error {
//	list, err := h.svc.GetTenantList()
//}
