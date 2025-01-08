package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/service/user"
	"sync"
)

const handlerName = "user"

type Handler struct {
	svc *user.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			svc: user.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.Group(h.name()).Use(mw.Guard)

	g.Get("/me", h.me)
}

// Me
// @Summary user profile.
// @Description user profile.
// @Tags me
// @Accept json
// @Produce json
// @Success 200 {object} model.Session
// @Failure default {object} mw.HTTPError
// @Router /user/me [get]
func (h *Handler) me(c *fiber.Ctx) error {
	if sess := mw.Sess(c); sess != nil {
		u, err := h.svc.GetByEmail(c.Context(), sess.User.Email, sess.User.Role)
		if err != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(u.Safe())
	}

	return fiber.ErrForbidden
}
