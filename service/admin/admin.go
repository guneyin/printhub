package admin

import (
	"context"
	"errors"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/repo/admin"
	"github.com/guneyin/printhub/service/auth"
	"github.com/guneyin/printhub/service/tenant"
	"github.com/guneyin/printhub/service/user"
	"github.com/guneyin/printhub/utils"
	"log/slog"
	"sync"
	"time"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	repo    *admin.Repo
	authSvc *auth.Service
	tenant  *tenant.Service
	userSvc *user.Service
}

func newService() *Service {
	s := &Service{
		repo:    admin.NewRepo(),
		authSvc: auth.GetService(),
		tenant:  tenant.GetService(),
		userSvc: user.GetService(),
	}
	s.boostrap()
	return s
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) boostrap() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if cnt, _ := s.repo.GetCount(ctx); cnt == 0 {
		pwd, err := utils.RandomString(10)
		if err != nil {
			slog.Error("random string error", "error:", err)
			return
		}

		u := &model.User{
			Role:     model.UserRoleAdmin,
			Email:    "admin@ph.com",
			Name:     "Admin",
			Password: pwd,
			Active:   true,
		}

		err = s.repo.Boostrap(ctx, u)
		if err != nil {
			slog.Error("boostrap admin user error", "error:", err)
			return
		}

		slog.Info("boostrap admin user", "user:", u.Email, "password:", pwd)
		return
	}

	return
}

func (s *Service) GetCount(ctx context.Context) (int64, error) {
	return s.repo.GetCount(ctx)
}

func (s *Service) GetTenantList(ctx context.Context) (model.TenantList, error) {
	return s.repo.GetTenantList(ctx)
}

func (s *Service) GetTenantById(ctx context.Context, filter string) (*model.Tenant, error) {
	qf := utils.NewQueryFilter(filter)
	id, ok := qf.Get("id")
	if !ok {
		return nil, errors.New("filter error")
	}

	return s.repo.GetTenantById(ctx, id.String())
}

func (s *Service) TenantCreate(ctx context.Context, t *model.Tenant) error {
	_, err := s.repo.CreateTenant(ctx, t)
	if err != nil {
		return err
	}

	return nil
}
