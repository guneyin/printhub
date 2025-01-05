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
			Expiration: time.Hour * 24 * 30,
			Storage: sqlite3.New(sqlite3.Config{
				Database: market.Get().Config.DbPath,
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
		slog.ErrorContext(c.Context(), "getSession", "error:", err.Error())
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
		slog.ErrorContext(c.Context(), "AuthorizeSession", "error:", err.Error())
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

func AdminGuard(c *fiber.Ctx) error {
	if Sess(c).IsValid(model.UserRoleAdmin) {
		return c.Next()
	}
	return fiber.ErrUnauthorized
}

func Sess(c *fiber.Ctx) *model.Session {
	s := getSession(c)
	sess := s.Get("session")
	if sess != nil {
		return sess.(*model.Session)
	}
	return &model.Session{}
}
