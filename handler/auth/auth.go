package auth

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/auth"
	"github.com/guneyin/printhub/service/user"
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

	g.Post("/register", h.RegisterUser)
	g.Post("/login", h.LoginUser)
	g.Get("/oauth/:provider", h.OAuthInit)
	g.Get("/oauth/:provider/complete", h.OAuthComplete)
	g.Get("/logout", h.LogoutUser)
	g.Get("/recover", h.RecoverPassword)
	g.Get("/verify", h.VerifyToken)
	g.Get("/change", h.ChangePassword)
	g.Get("/validate", h.ValidateUser)
}

// @Router /auth/register [post].
func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	ur, err := model.NewAuthUserRequest(c.Body())
	if err != nil {
		return err
	}
	u := &model.User{
		Role:     role,
		Email:    ur.Email,
		Password: ur.Password,
	}

	err = h.svc.RegisterUser(c.Context(), u)
	if err != nil {
		return err
	}

	return c.JSON(u.Safe())
}

// @Router /auth/login [post].
func (h *Handler) LoginUser(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	ur, err := model.NewAuthUserRequest(c.Body())
	if err != nil {
		return err
	}

	sess, err := h.svc.LoginUser(c.Context(), role, ur.Email, ur.Password)
	if err != nil {
		return err
	}

	err = mw.AuthorizeSession(c, sess)
	if err != nil {
		return err
	}

	return c.JSON(sess)
}

// @Router /auth/oauth/{provider} [get].
func (h *Handler) OAuthInit(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	u, err := h.svc.InitOAuth(
		c.Params("provider"),
		role,
		// c.Query("callback"),
		c.QueryBool("force"))
	if err != nil {
		return err
	}

	return c.Redirect(u, http.StatusFound)
}

func (h *Handler) OAuthComplete(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	sess, err := h.svc.CompleteOAuth(c.Context(), role, c.Params("provider"), c.Query("code"))
	if err != nil {
		return err
	}

	err = mw.AuthorizeSession(c, sess)
	if err != nil {
		return err
	}

	return c.JSON(sess)
}

// @Router /auth/logout [post].
func (h *Handler) LogoutUser(c *fiber.Ctx) error {
	return mw.InvalidateSession(c)
}

func (h *Handler) RecoverPassword(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	email := c.Query("email")

	h.svc.RecoverPassword(c.Context(), email, role)

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) VerifyToken(c *fiber.Ctx) error {
	token := c.Query("token")
	u, err := h.svc.VerifyToken(c.Context(), token)
	if err != nil {
		return err
	}
	return c.JSON(u.Safe())
}

func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	token := c.Query("token")
	password := c.Query("password")

	return h.svc.ChangePassword(c.Context(), token, password)
}

func (h *Handler) ValidateUser(c *fiber.Ctx) error {
	token := c.Query("token")
	u, err := h.svc.ValidateUser(c.Context(), token)
	if err != nil {
		return err
	}
	return c.JSON(u.Safe())
}
