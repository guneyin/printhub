package model

import (
	"encoding/gob"
)

func init() {
	gob.Register(&Session{})
}

type Session struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	User     User   `json:"user"`
}

// todo: s.user.role boş
func (s *Session) IsValid(role UserRole) bool {
	return s.User.Role == role
}
