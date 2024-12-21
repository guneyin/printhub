package user

import (
	"context"
	"errors"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/repo/user"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	repo *user.Repo
}

func newService() *Service {
	return &Service{repo: user.NewRepo()}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) InitUser(ctx context.Context, u *model.User) (*model.User, error) {
	found, _ := s.repo.GetByEmail(ctx, u.Email, u.Role)

	switch u.Role {
	case model.UserRoleClient:
		return s.initClientUser(ctx, found, u)
	case model.UserRoleTenant:
		return s.initTenantUser(ctx, found, u)
	default:
		return nil, errors.New("invalid user role")
	}
}

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	return s.repo.GetByUUID(ctx, uuid)
}

func (s *Service) GetByEmail(ctx context.Context, email string, role model.UserRole) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email, role)
}

func (s *Service) initClientUser(ctx context.Context, found, u *model.User) (*model.User, error) {
	if found != nil {
		return found, nil
	}
	err := s.repo.Create(ctx, u)
	return u, err
}

func (s *Service) initTenantUser(ctx context.Context, found, u *model.User) (*model.User, error) {
	if found == nil {
		return nil, errors.New("tenant user not found! contact to admin")
	}
	if found.IsActivated() {
		return found, nil
	}
	err := s.repo.Create(ctx, u)
	return u, err
}
