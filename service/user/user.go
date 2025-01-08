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

func (s *Service) Create(ctx context.Context, u *model.User) (*model.User, error) {
	found, _ := s.repo.GetByEmail(ctx, u.Email, u.Role)
	if found != nil {
		return nil, errors.New("user already exists")
	}

	return s.repo.Create(ctx, u)
}

func (s *Service) Update(ctx context.Context, uuid string, u *model.User) (*model.User, error) {
	found, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	u.ID = found.ID

	return s.repo.Update(ctx, u)
}

func (s *Service) InitUser(ctx context.Context, u *model.User) (*model.User, error) {
	switch u.Role {
	case model.UserRoleAdmin:
		return nil, errors.New("admin user not allowed")
	case model.UserRoleTenant:
		return s.initTenantUser(ctx, u)
	case model.UserRoleClient:
		return s.initClientUser(ctx, u)
	default:
		return nil, errors.New("invalid user role")
	}
}

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	return s.repo.GetByUUID(ctx, uuid)
}

func (s *Service) GetByEmail(ctx context.Context, email string, role model.UserRole) (*model.User, error) {
	u, err := s.repo.GetByEmail(ctx, email, role)
	if err != nil {
		return nil, err
	}

	if !u.IsActivated() {
		return nil, errors.New("user is not active")
	}

	return u, nil
}

func (s *Service) GetByEmailWithoutRestriction(ctx context.Context, email string, role model.UserRole) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email, role)
}

func (s *Service) initUser(ctx context.Context, email string, role model.UserRole) (*model.User, error) {
	u, err := s.repo.GetByEmail(ctx, email, role)
	if err != nil {
		return s.repo.Create(ctx, &model.User{
			Email:  email,
			Role:   role,
			Active: true,
		})
	}

	if !u.IsActivated() {
		u.Active = true

		return s.repo.Update(ctx, u)
	}

	return u, nil
}

func (s *Service) initClientUser(ctx context.Context, u *model.User) (*model.User, error) {
	return s.initUser(ctx, u.Email, u.Role)
}

func (s *Service) initTenantUser(ctx context.Context, u *model.User) (*model.User, error) {
	found, _ := s.repo.GetByEmail(ctx, u.Email, u.Role)
	if found == nil {
		return nil, errors.New("tenant user not found! contact to admin")
	}

	return s.initUser(ctx, u.Email, u.Role)
}

func (s *Service) Validate(ctx context.Context, u *model.User) error {
	u.Active = true
	_, err := s.repo.Update(ctx, u)
	return err
}
