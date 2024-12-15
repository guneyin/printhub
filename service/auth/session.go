package auth

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&Token{})
	gob.Register(&OAuthUser{})
}

type Session struct {
	Token Token     `json:"token"`
	User  OAuthUser `json:"user,omitempty"`
}

type Token struct {
	Provider     string    `json:"provider,omitempty"`
	AccessToken  string    `json:"accessToken,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
	IDToken      string    `json:"idToken,omitempty"`
}

type OAuthUser struct {
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Gender    string `json:"gender,omitempty"`
	UserID    string `json:"userID,omitempty"`
	AvatarURL string `json:"avatarURL,omitempty"`
	Location  string `json:"location,omitempty"`
	Link      string `json:"link,omitempty"`
}
