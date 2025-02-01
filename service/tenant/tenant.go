package tenant

import (
	"context"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/repo/tenant"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	repo *tenant.Repo
}

func newService() *Service {
	return &Service{repo: tenant.NewRepo()}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) Save(ctx context.Context, t *model.Tenant) error {
	_, err := s.repo.Save(ctx, t)
	return err
}

func (s *Service) Search(ctx context.Context, filter model.QueryFilter) (model.TenantList, error) {
	return s.repo.Search(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	filter := new(model.TenantFilter)
	filter.ID = id

	list, err := s.repo.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &list[0], nil
}

func (s *Service) Delete(ctx context.Context, t *model.Tenant) error {
	return nil
}

func (s *Service) AddUser(ctx context.Context, t *model.Tenant, u *model.User) error {
	return s.repo.AddUser(ctx, t, u)
}

func (s *Service) GetUserList(ctx context.Context, id string) (model.UserList, error) {

}

//func (s *Service) GetConfig(ctx context.Context, key string) (*model.ConfigList, error) {
//	return s.config.Get(ctx, key)
//}
//
//func (s *Service) SetConfig(ctx context.Context, list *model.ConfigList) error {
//	return s.config.Set(ctx, list)
//}
//
//func (s *Service) DeleteConfig(ctx context.Context, key string) error {
//	return s.config.Delete(ctx, key)
//}

//func (s *Service) DiskAuth(ctx context.Context, provider string) (string, error) {
//	disk, err := s.getDisk(ctx, provider)
//	if err != nil {
//		return "", err
//	}
//
//	return disk.InitAuth(), nil
//}
//
//func (s *Service) DiskAuthCallback(ctx context.Context, provider, code string) (*oauth2.Token, error) {
//	disk, err := s.getDisk(ctx, provider)
//	if err != nil {
//		return nil, err
//	}
//
//	return disk.VerifyAuth(ctx, code)
//}
//
//func (s *Service) getDisk(ctx context.Context, provider string) (disgo.Provider, error) {
//	if provider == "" {
//		return nil, fmt.Errorf("invalid provided")
//	}
//
//	key := fmt.Sprintf("%s:%s", provider, "config")
//	config, err := s.GetConfig(ctx, key)
//	if err != nil {
//		return nil, err
//	}
//	if len(*config) == 0 {
//		return nil, fmt.Errorf("no config found")
//	}
//	configData := (*config)[0].Value
//
//	key = fmt.Sprintf("%s:%s", provider, "token")
//	token, err := s.GetConfig(ctx, key)
//	if err != nil {
//		return nil, err
//	}
//
//	var tokenData []byte
//	if len(*token) != 0 {
//		tokenData = (*token)[0].JSON()
//	}
//
//	disk, err := disgo.New(ctx, provider, []byte(configData), tokenData)
//	if err != nil {
//		return nil, err
//	}
//	return disk, nil
//}
