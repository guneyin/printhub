package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/handler/mw"
	"sync"
)

const handlerName = "profile"

type Handler struct{}

var (
	once    sync.Once
	handler *Handler
)

func InitHandler(r fiber.Router) {
	once.Do(func() {
		handler = &Handler{}
		handler.setRoutes(r)
	})
}

func (h *Handler) name() string {
	return handlerName
}

func (h *Handler) setRoutes(r fiber.Router) {
	g := r.
		Group(h.name()).
		Use(mw.Protected())

	g.Get("/", h.profile)
	g.Get("/session", h.session)
}

func (h *Handler) profile(c *fiber.Ctx) error {
	if user := mw.User(c); user != nil {
		return c.JSON(user)
	}
	return fiber.ErrForbidden
}

func (h *Handler) session(c *fiber.Ctx) error {
	sess := mw.GetSession(c)
	keys := sess.Keys()

	res := map[string]interface{}{}
	for _, key := range keys {
		res[key] = sess.Get(key)
	}

	return c.JSON(&res)
}
