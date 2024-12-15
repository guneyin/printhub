package mw

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/guneyin/printhub/service/auth"
	"log/slog"
)

//var store = session.New(session.Config{Storage: sqlite3.New(sqlite3.Config{
//	Database: market.Get().Config.DbPath,
//	Table:    "sessions",
//})})

var store = session.New()

func GetSession(c *fiber.Ctx) *session.Session {
	ss, err := store.Get(c)
	if err != nil {
		slog.ErrorContext(c.Context(), "getSession", "error:", err.Error())
		return &session.Session{}
	}
	return ss
}

type HTTPError struct {
	Error string `json:"error"`
}

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if User(c) != nil {
			return c.Next()
		}
		return fiber.ErrUnauthorized
	}
}

func User(c *fiber.Ctx) *auth.OAuthUser {
	if user := GetSession(c).Get("user"); user != nil {
		return user.(*auth.OAuthUser)
	}
	return nil
}
