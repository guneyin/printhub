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

// RegisterUser Register
// @Summary Register client user.
// @Description Register client user.
// @Tags Auth Register
// @Accept json
// @Produce json
// @Param role query string true "role" Enums(admin, tenant, client)
// @Param tenant body model.AuthUserRequest true "login info"
// @Failure default {object} mw.HTTPError
// @Router /auth/register [post]
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

// LoginUser Login
// @Summary login.
// @Description login.
// @Tags login
// @Accept json
// @Produce json
// @Param role query string true "role" Enums(admin, tenant, client)
// @Param tenant body model.AuthUserRequest true "login info"
// @Success 200 {object} model.Session
// @Failure default {object} mw.HTTPError
// @Router /auth/login [post]
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

// OAuthInit Init
// @Summary Init auth.
// @Description Start OAuth2 authorization.
// @Tags Auth OAuthInit
// @Accept json
// @Produce json
// @Param provider path string true "provider"
// @Param role query string true "role"
// @Param callback query string true "callback url"
// @Param force query bool false "force"
// @Failure default {object} mw.HTTPError
// @Router /auth/oauth/{provider} [get]
func (h *Handler) OAuthInit(c *fiber.Ctx) error {
	role, err := model.NewUserRole(c.Query("role"))
	if err != nil {
		return err
	}

	u, err := h.svc.InitOAuth(
		c.Params("provider"),
		role,
		c.Query("callback"),
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

// LogoutUser Logout
// @Summary Logout.
// @Description Logout.
// @Tags Logout
// @Accept json
// @Produce json
// @Success 200 {object} model.Session
// @Failure default {object} mw.HTTPError
// @Router /auth/logout [post]
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
