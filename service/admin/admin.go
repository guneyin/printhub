package admin

import (
	"context"
	"sync"
	"time"

	"github.com/guneyin/printhub/market"

	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/repo/admin"
	"github.com/guneyin/printhub/service/auth"
	"github.com/guneyin/printhub/service/tenant"
	"github.com/guneyin/printhub/service/user"
	"github.com/guneyin/printhub/utils"
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
			market.Log().Error("random string error", "error:", err)
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
			market.Log().Error("boostrap admin user error", "error:", err)
			return
		}

		market.Log().Info("boostrap admin user", "user:", u.Email, "password:", pwd)
		return
	}
}

func (s *Service) GetCount(ctx context.Context) (int64, error) {
	return s.repo.GetCount(ctx)
}

func (s *Service) SaveTenant(ctx context.Context, t *model.Tenant) error {
	return s.tenant.Save(ctx, t)
}

func (s *Service) SearchTenant(ctx context.Context, filter model.QueryFilter) (model.TenantList, error) {
	return s.tenant.Search(ctx, filter)
}

func (s *Service) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	return s.tenant.GetByID(ctx, id)
}

func (s *Service) DeleteTenant(ctx context.Context, id string) error {
	return s.tenant.Delete(ctx, id)
}

func (s *Service) AddTenantUser(ctx context.Context, tenantID string, user *model.User) error {
	return s.tenant.AddUser(ctx, tenantID, user)
}

func (s *Service) GetTenantUserList(ctx context.Context, id string) (model.UserList, error) {
	return s.tenant.GetUserList(ctx, id)
}
