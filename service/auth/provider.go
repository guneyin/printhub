package auth

import (
	"context"
	"fmt"
	"github.com/guneyin/printhub/model"
)

const (
	OAuthProviderGoogle = "google"
)

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
