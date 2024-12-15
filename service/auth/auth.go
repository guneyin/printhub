package auth

import (
	"context"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct{}

func newService() *Service {
	return &Service{}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) InitOAuth(provider string, force bool) (string, error) {
	p, err := NewProvider(provider)
	if err != nil {
		return "", err
	}

	return p.InitOAuth(force)
}

func (s *Service) CompleteOAuth(ctx context.Context, provider string, code string) (*Session, error) {
	p, err := NewProvider(provider)
	if err != nil {
		return nil, err
	}

	return p.CompleteOAuth(ctx, code)
}
