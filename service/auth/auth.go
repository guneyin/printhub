package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/guneyin/printhub/mail"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/user"
	"github.com/guneyin/printhub/utils"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"strings"
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

func (s *Service) RegisterUser(ctx context.Context, u *model.User) error {
	if u.Role != model.UserRoleClient {
		return fmt.Errorf("%s user not allowed", u.Role)
	}

	if err := u.Validate(); err != nil {
		return err
	}

	created, err := s.userSvc.Create(ctx, u)
	if err != nil {
		return err
	}

	token, err := generateToken(created.UUID)
	if err != nil {
		slog.Warn(err.Error())
		return nil
	}

	rp := mail.NewVerifyUserEmail(token)
	err = rp.Send(created.Email, "Hesabınızı doğrulayın")
	if err != nil {
		slog.Warn(err.Error())
	}

	return nil
}

func (s *Service) InitOAuth(provider string, role model.UserRole, cbUrl string, force bool) (string, error) {
	p, err := NewProvider(provider)
	if err != nil {
		return "", err
	}

	return p.InitOAuth(role, cbUrl, force)
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

	u := oauth.ToUser(role)

	created, err := s.userSvc.InitUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return s.createSession("oauth", created)
}

func (s *Service) LoginUser(ctx context.Context, role model.UserRole, email, password string) (*model.Session, error) {
	u, err := s.userSvc.GetByEmail(ctx, email, role)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid auth credentials")
	}

	return s.createSession("basic", u)
}

func (s *Service) createSession(provider string, u *model.User) (*model.Session, error) {
	return &model.Session{
		Provider: provider,
		User:     *u.Safe(),
	}, nil
}

func (s *Service) RecoverPassword(ctx context.Context, email string, role model.UserRole) {
	u, err := s.userSvc.GetByEmail(ctx, email, role)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	if !u.IsActivated() {
		return
	}

	token, err := generateToken(u.UUID)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	rp := mail.NewRecoverPasswordEmail(token)
	err = rp.Send(email, "Parola Sıfırlama")
	if err != nil {
		slog.Warn(err.Error())
	}
}

func (s *Service) VerifyToken(ctx context.Context, token string) (*model.User, error) {
	uuid, err := verifyToken(token)
	if err != nil {
		return nil, err
	}

	return s.userSvc.GetByUUID(ctx, uuid)
}

func (s *Service) ChangePassword(ctx context.Context, token, password string) error {
	uuid, err := verifyToken(token)
	if err != nil {
		return err
	}

	u := &model.User{Password: password}
	_, err = s.userSvc.Update(ctx, uuid, u)
	return err
}

func (s *Service) ValidateUser(ctx context.Context, token string) (*model.User, error) {
	u, err := s.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if u.IsActivated() {
		return nil, errors.New("user is already validated")
	}

	err = s.userSvc.Validate(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func generateToken(uid string) (string, error) {
	cfg := market.Get().Config
	tokenStr := fmt.Sprintf("%s:%d", uid, time.Now().Unix())
	return utils.Encrypt(tokenStr, []byte(cfg.AuthSecret))
}

func verifyToken(hash string) (string, error) {
	cfg := market.Get().Config
	token, err := utils.Decrypt(hash, []byte(cfg.AuthSecret))
	if err != nil {
		return "", err
	}

	split := strings.Split(token, ":")
	if len(split) != 2 {
		return "", errors.New("invalid token")
	}

	uuid, ts := split[0], split[1]
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return "", errors.New("invalid token")
	}

	since := time.Since(time.Unix(timestamp, 0))
	if since > time.Minute*10 {
		return "", errors.New("token expired")
	}

	return uuid, nil
}
