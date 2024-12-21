package auth

import (
	"context"
	"fmt"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gapi "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"net/url"
)

var (
	googleAuthScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.appdata",
		"https://www.googleapis.com/auth/drive.file",
		"https://www.googleapis.com/auth/drive.metadata",
		"https://www.googleapis.com/auth/drive.metadata.readonly",
		"https://www.googleapis.com/auth/drive.photos.readonly",
		"https://www.googleapis.com/auth/drive.readonly",
	}
)

type googleProvider struct{}

func newGoogleProvider() *googleProvider {
	return &googleProvider{}
}

func (gp *googleProvider) config() *oauth2.Config {
	cfg := market.Get().Config
	u := fmt.Sprintf("%s/api/auth/google/callback", cfg.ApiBaseUrl)

	return &oauth2.Config{
		RedirectURL:  u,
		ClientID:     cfg.GoogleClientId,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       googleAuthScopes,
		Endpoint:     google.Endpoint,
	}
}

func (gp *googleProvider) InitOAuth(role model.UserRole, force bool) (string, error) {
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOnline,
	}
	if force {
		opts = append(opts, oauth2.ApprovalForce)
	}

	u := gp.config().AuthCodeURL(string(role), opts...)
	au, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return au.String(), nil
}

func (gp *googleProvider) CompleteOAuth(ctx context.Context, code string) (*OAuthUser, error) {
	token, err := gp.config().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	ts := gp.config().TokenSource(ctx, token)
	svc, err := gapi.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}

	user, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	oauth := &OAuthUser{
		Email:     user.Email,
		Name:      user.Name,
		LastName:  user.FamilyName,
		Gender:    user.Gender,
		UserID:    user.Id,
		AvatarURL: user.Picture,
		Location:  user.Locale,
		Link:      user.Link,
	}

	return oauth, nil
}
