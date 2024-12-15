package config

import (
	"context"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/repo/config"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	repo *config.Repo
}

func newService() *Service {
	return &Service{repo: config.NewConfigRepo()}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) Get(ctx context.Context, id, module, key string) (*model.Config, error) {
	return s.repo.Get(ctx, id, module, key)
}

func (s *Service) Set(ctx context.Context, list *model.ConfigList) error {
	return s.repo.Set(ctx, list)
}

func (s *Service) Delete(ctx context.Context, id, module, key string) error {
	return s.repo.Delete(ctx, id, module, key)
}
