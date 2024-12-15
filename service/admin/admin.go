package admin

import (
	"context"
	"errors"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/repo/admin"
	"github.com/guneyin/printhub/repo/tenant"
	"github.com/guneyin/printhub/repo/user"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	repo   *admin.Repo
	tenant *tenant.Repo
	user   *user.Repo
}

func newService() *Service {
	return &Service{
		repo:   admin.NewRepo(),
		tenant: tenant.NewRepo(),
		user:   user.NewRepo(),
	}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) TenantCreate(ctx context.Context, t *model.Tenant) error {
	err := s.tenant.Create(ctx, t)
	if err != nil {
		return err
	}

	u := &model.User{
		Role:  model.UserRoleTenant,
		Email: t.Email,
	}
	err = s.user.Create(ctx, u)
	if err != nil {
		return errors.New("tenant created but user create failed")
	}

	err = s.tenant.AddUser(ctx, t, u)
	if err != nil {
		return errors.New("tenant created but add user failed")
	}

	return nil
}
