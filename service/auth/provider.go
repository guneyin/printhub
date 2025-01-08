package auth

import (
	"context"
	"fmt"
	"github.com/guneyin/printhub/model"
)

const (
	OAuthProviderGoogle = "google"
)

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

type Provider interface {
	InitOAuth(role model.UserRole, cbUrl string, force bool) (string, error)
	CompleteOAuth(ctx context.Context, code string) (*OAuthUser, error)
}

func NewProvider(provider string) (Provider, error) {
	switch provider {
	case OAuthProviderGoogle:
		return newGoogleProvider(), nil
	}
	return nil, fmt.Errorf("unknown provider: %s", provider)
}

func (o *OAuthUser) ToUser(role model.UserRole) *model.User {
	return &model.User{
		Role:      role,
		Email:     o.Email,
		Name:      o.Email,
		Password:  "",
		AvatarURL: o.AvatarURL,
	}
}
