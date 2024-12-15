package auth

import (
	"context"
	"fmt"
)

const (
	OAuthProviderGoogle = "google"
)

type Provider interface {
	InitOAuth(force bool) (string, error)
	CompleteOAuth(ctx context.Context, code string) (*Session, error)
}

func NewProvider(provider string) (Provider, error) {
	switch provider {
	case OAuthProviderGoogle:
		return newGoogleProvider(), nil
	}
	return nil, fmt.Errorf("unknown provider: %s", provider)
}
