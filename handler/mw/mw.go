package mw

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"log/slog"
	"sync"
	"time"
)

var (
	onceStore sync.Once
	ss        *session.Store
)

func store() *session.Store {
	onceStore.Do(func() {
		ss = session.New(session.Config{
			Expiration: time.Hour * 72,
			Storage: sqlite3.New(sqlite3.Config{
				Database: market.Get().Config.DbPath,
				Table:    "sessions",
			}),
			KeyLookup:         "cookie:session_id",
			CookieDomain:      "",
			CookiePath:        "",
			CookieSecure:      true,
			CookieHTTPOnly:    true,
			CookieSameSite:    "Strict",
			CookieSessionOnly: false,
			KeyGenerator:      utils.UUIDv4,
			CookieName:        "",
		})
	})
	return ss
}

func getSession(c *fiber.Ctx) *session.Session {
	s, err := store().Get(c)
	if err != nil {
		slog.ErrorContext(c.Context(), "getSession", "error:", err.Error())
		return &session.Session{}
	}
	return s
}

func AuthorizeSession(c *fiber.Ctx, sess *model.Session) {
	s := getSession(c)
	s.Set("session", sess)
	_ = s.Save()
}

type HTTPError struct {
	Error string `json:"error"`
}

func AdminGuard(c *fiber.Ctx) error {
	c.Locals("role", model.UserRoleAdmin)
	return protected(c)
}

func TenantGuard(c *fiber.Ctx) error {
	c.Locals("role", model.UserRoleTenant)
	return protected(c)
}

func ClientGuard(c *fiber.Ctx) error {
	c.Locals("role", model.UserRoleClient)
	return protected(c)
}

func protected(c *fiber.Ctx) error {
	role := c.Locals("role", model.UserRoleAdmin).(model.UserRole)
	if Sess(c).IsValid(role) {
		return c.Next()
	}
	return fiber.ErrUnauthorized

}

func Sess(c *fiber.Ctx) *model.Session {
	if s := getSession(c).Get("session"); s != nil {
		return s.(*model.Session)
	}
	return &model.Session{}
}
