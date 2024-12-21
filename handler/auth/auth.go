package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/auth"
	"github.com/guneyin/printhub/service/user"
	"net/http"
	"sync"
)

const handlerName = "auth"

type Handler struct {
	svc     *auth.Service
	userSvc *user.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc:     auth.GetService(),
			userSvc: user.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.Group(h.name())

	g.Get("/oauth/:provider", h.InitOAuth)
	g.Get("/oauth/:provider/CompleteOAuth", h.CompleteOAuth)
	g.Post("/login", h.BasicAuth)
}

// InitOAuth Init
// @Summary Init auth.
// @Description Start OAuth2 authorization.
// @Tags Auth InitOAuth
// @Accept json
// @Produce json
// @Param provider path string true "provider"
// @Param role query string true "role"
// @Param force query bool false "force"
// @Failure default {object} mw.HTTPError
// @Router /auth/oauth/{provider} [get]
func (h *Handler) InitOAuth(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	u, err := h.svc.InitOAuth(
		c.Params("provider"),
		role,
		c.QueryBool("force"))
	if err != nil {
		return err
	}

	return c.Redirect(u, http.StatusFound)
}

func (h *Handler) CompleteOAuth(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("state"))
	if err != nil {
		return err
	}

	sess, err := h.svc.CompleteOAuth(c.Context(), role, c.Params("provider"), c.Query("code"))
	if err != nil {
		return err
	}

	mw.AuthorizeSession(c, sess)

	return c.SendStatus(fiber.StatusOK)
}

// BasicAuth
// @Summary login.
// @Description login.
// @Tags login
// @Accept json
// @Produce json
// @Param role query string true "role" Enums(admin, tenant, client)
// @Param tenant body model.AuthLoginRequest true "login info"
// @Failure default {object} mw.HTTPError
// @Router /auth/login [post]
func (h *Handler) BasicAuth(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	req, err := model.NewAuthLoginRequest(c.Body())
	if err != nil {
		return err
	}

	sess, err := h.svc.BasicAuth(c.Context(), role, req.Email, req.Password)
	if err != nil {
		return err
	}

	mw.AuthorizeSession(c, sess)

	return c.SendStatus(fiber.StatusOK)
}
