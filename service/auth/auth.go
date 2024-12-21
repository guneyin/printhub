package auth

import (
	"context"
	"errors"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/user"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	userSvc *user.Service
}

func newService() *Service {
	return &Service{userSvc: user.GetService()}
}

func GetService() *Service {
	once.Do(func() {
		service = newService()
	})
	return service
}

func (s *Service) InitOAuth(provider string, role model.UserRole, force bool) (string, error) {
	p, err := NewProvider(provider)
	if err != nil {
		return "", err
	}

	return p.InitOAuth(role, force)
}

func (s *Service) CompleteOAuth(ctx context.Context, role model.UserRole, provider, code string) (*model.Session, error) {
	p, err := NewProvider(provider)
	if err != nil {
		return nil, err
	}

	oauth, err := p.CompleteOAuth(ctx, code)
	if err != nil {
		return nil, err
	}

	u := oauth.ToUser()
	u.Role = role
	return s.getSession("oauth", u)
}

func (s *Service) BasicAuth(ctx context.Context, role model.UserRole, email, password string) (*model.Session, error) {
	u, err := s.userSvc.GetByEmail(ctx, email, role)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid auth credentials")
	}

	return s.getSession("basic", u)
}

func (s *Service) getSession(provider string, u *model.User) (*model.Session, error) {
	return &model.Session{
		Provider:  provider,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		UserId:    u.UUID,
		UserEmail: u.Email,
		UserRole:  u.Role,
	}, nil
}
