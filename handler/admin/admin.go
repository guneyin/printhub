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
	tenant.Get("/list", h.getTenantList)
	tenant.Get("/", h.getTenantById)
	tenant.Post("/", h.tenantCreate)
	tenant.Post("/user", h.tenantUserCreate)
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
	list, err := h.svc.GetTenantList(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(list)
}

func (h *Handler) getTenantById(c *fiber.Ctx) error {
	tenant, err := h.svc.GetTenantById(c.Context(), c.Query("filter"))
	if err != nil {
		return err
	}
	return c.JSON(tenant)
}

// CreateTenant
// @Summary tenant create.
// @Description Create a new tenant.
// @Tags tenant create
// @Accept json
// @Produce json
// @Param tenant body model.Tenant true "tenant"
// @Failure default {object} mw.HTTPError
// @Router /admin/tenant [post]
func (h *Handler) tenantCreate(c *fiber.Ctx) error {
	body := c.Body()
	bodyStr := string(body)
	tenant, err := model.NewTenant([]byte(bodyStr))
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
