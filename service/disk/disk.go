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
	key := fmt.Sprintf("%s:%s", provider, utils.ConfigParamOAuth.String())
	oauth2, err := s.cfr.Get(ctx, identifier, utils.ConfigParamDisk.String(), key)
	if err != nil {
		return nil, err
	}

	key = fmt.Sprintf("%s:%s", provider, utils.ConfigParamToken.String())
	token, _ := s.cfr.Get(ctx, identifier, utils.ConfigParamDisk.String(), key)

	disk, err := disgo.New(ctx, disgo.ProviderType(provider), oauth2.JSON(), token.JSON())
	if err != nil {
		return nil, err
	}
	return disk, nil
}

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
