package admin

import (
	"context"
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
	repo   *admin.Repo
	auth   *auth.Service
	tenant *tenant.Service
	user   *user.Service
}

func newService() *Service {
	s := &Service{
		repo:   admin.NewRepo(),
		auth:   auth.GetService(),
		tenant: tenant.GetService(),
		user:   user.GetService(),
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

func (s *Service) SearchTenant(ctx context.Context, filter model.QueryFilter) (model.TenantList, error) {
	return s.tenant.Search(ctx, filter)
}

func (s *Service) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	filter := new(model.TenantFilter)
	filter.ID = id

	list, err := s.tenant.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &list[0], nil
}

func (s *Service) SaveTenant(ctx context.Context, t *model.Tenant) error {
	return s.tenant.Save(ctx, t)
}
