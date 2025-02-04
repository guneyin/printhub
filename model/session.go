package model

import (
	"slices"

	"github.com/google/uuid"
)

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
