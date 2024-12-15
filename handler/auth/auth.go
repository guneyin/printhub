package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/service/auth"
	"net/http"
	"strings"
	"sync"
)

const handlerName = "auth"

type Handler struct {
	svc *auth.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: auth.GetService(),
		}
		handler.setRoutes(r)
	})
}
func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.Group(h.name())

	g.Get("/:provider", h.init)
	g.Get("/:provider/callback", h.callback)
}

// Init
// @Summary Init auth.
// @Description Start OAuth2 authorization.
// @Tags Auth init
// @Accept json
// @Produce json
// @Param provider path string true "provider"
// @Failure all {object} mw.HTTPError{}
// @Router /auth/{provider} [get]
func (h *Handler) init(c *fiber.Ctx) error {
	u, err := h.svc.InitOAuth(c.Params("provider"), c.QueryBool("force"))
	if err != nil {
		return err
	}

	if role := strings.TrimSpace(c.Query("role")); role != "" {
		sess := mw.GetSession(c)
		sess.Set("role", role)
		_ = sess.Save()
	}

	return c.Redirect(u, http.StatusFound)
}

func (h *Handler) callback(c *fiber.Ctx) error {
	sessionData, err := h.svc.CompleteOAuth(c.Context(), c.Params("provider"), c.Query("code"))
	if err != nil {
		return err
	}

	//todo: keep working here
	sess := mw.GetSession(c)
	//role := sess.Get("role")

	if err = sess.Reset(); err != nil {
		return err
	}

	sess.Set("user", sessionData.User)

	if err = sess.Save(); err != nil {
		return err
	}

	return c.JSON(sessionData.User)
}
