package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"github.com/guneyin/printhub/service/user"
	"sync"
)

const handlerName = "profile"

type Handler struct {
	userSvc *user.Service
}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{
			userSvc: user.GetService(),
		}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.
		Group(h.name()).
		Use(mw.AdminGuard)

	g.Get("/", h.profile)
}

func (h *Handler) profile(c *fiber.Ctx) error {
	if sess := mw.Sess(c); sess != nil {
		u, err := h.userSvc.GetByEmail(c.Context(), sess.UserEmail, sess.UserRole)
		if err != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(u.Safe())
	}

	return fiber.ErrForbidden
}
