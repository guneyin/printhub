package model

import (
	"encoding/gob"
	"github.com/google/uuid"
	"slices"
)

func init() {
	gob.Register(&Session{})
}

type Session struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	User     User   `json:"user"`
}

func (s *Session) IsAuthorized(role ...UserRole) bool {
	_, err := uuid.Parse(s.ID)
	if err != nil {
		return false
	}

	if role == nil {
		return true
	}
	return slices.Contains(role, s.User.Role)
}
