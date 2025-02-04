package mw

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
)

var (
	onceStore sync.Once
	ss        *session.Store
)

func store() *session.Store {
	onceStore.Do(func() {
		ss = session.New(session.Config{
			Expiration: time.Hour * 24 * 30,
			Storage: sqlite3.New(sqlite3.Config{
				Database: market.Get().Config.DBPath,
				Table:    "sessions",
			}),
			KeyLookup:         "cookie:session_id",
			CookieDomain:      "",
			CookiePath:        "/",
			CookieSecure:      true,
			CookieHTTPOnly:    true,
			CookieSameSite:    "Strict",
			CookieSessionOnly: false,
			KeyGenerator:      utils.UUIDv4,
		})
	})
	return ss
}

func getSession(c *fiber.Ctx) *session.Session {
	s, err := store().Get(c)
	if err != nil {
		market.Log().ErrorContext(c.Context(), "getSession", "error:", err.Error())
		return &session.Session{}
	}
	return s
}

func AuthorizeSession(c *fiber.Ctx, sess *model.Session) error {
	s := getSession(c)
	sess.ID = s.ID()
	s.Set("session", sess)
	err := s.Save()
	if err != nil {
		market.Log().ErrorContext(c.Context(), "AuthorizeSession", "error:", err.Error())
		return err
	}

	return nil
}

func InvalidateSession(c *fiber.Ctx) error {
	s := getSession(c)
	return s.Destroy()
}

type HTTPError struct {
	Error string `json:"error"`
}

func Guard(c *fiber.Ctx) error {
	if Sess(c).IsAuthorized() {
		return c.Next()
	}
	return fiber.ErrUnauthorized
}

func AdminGuard(c *fiber.Ctx) error {
	if Sess(c).IsAuthorized(model.UserRoleAdmin) {
		return c.Next()
	}
	return fiber.ErrUnauthorized
}

func TenantGuard(c *fiber.Ctx) error {
	if Sess(c).IsAuthorized(model.UserRoleTenant) {
		return c.Next()
	}
	return fiber.ErrUnauthorized
}

func ClientGuard(c *fiber.Ctx) error {
	if Sess(c).IsAuthorized(model.UserRoleClient) {
		return c.Next()
	}
	return fiber.ErrUnauthorized
}

func Sess(c *fiber.Ctx) *model.Session {
	sess := getSession(c).Get("session")

	if s := sess.(*model.Session); s != nil {
		return s
	}

	return &model.Session{}
}
