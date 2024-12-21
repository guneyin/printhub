package model

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&Session{})
}

type Session struct {
	Provider     string    `json:"provider"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	UserId       string    `json:"userId"`
	UserEmail    string    `json:"userEmail"`
	UserRole     UserRole  `json:"userRole"`
}

func (s *Session) IsValid(role UserRole) bool {
	return s.UserRole == role
}
