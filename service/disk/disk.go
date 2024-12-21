package disk

import (
	"context"
	"fmt"
	"github.com/guneyin/disgo"
	"github.com/guneyin/printhub/repo/config"
	"github.com/guneyin/printhub/utils"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	cfr *config.Repo
	//repo     *disk.Repo
}

func newService() *Service {
	return &Service{
		cfr: config.NewConfigRepo(),
	}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) disk(ctx context.Context, identifier, provider string) (disgo.Provider, error) {
	pt, err := disgo.NewProviderType(provider)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s:%s", provider, utils.ConfigParamOAuth.String())
	oauth2, err := s.cfr.Get(ctx, identifier, utils.ConfigParamDisk.String(), key)
	if err != nil {
		return nil, err
	}

	key = fmt.Sprintf("%s:%s", provider, utils.ConfigParamToken.String())
	token, _ := s.cfr.Get(ctx, identifier, utils.ConfigParamDisk.String(), key)

	disk, err := disgo.New(ctx, pt, oauth2.JSON(), token.JSON())
	if err != nil {
		return nil, err
	}
	return disk, nil
}
