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

	//u := cbUrl
	//if cbUrl == "" {
	//u := fmt.Sprintf("%s/api/auth/google/callback", cfg.ApiBaseUrl)
	//}
	u := "http://localhost:5173/oauth/google/callback"
	return &oauth2.Config{
		RedirectURL:  u,
		ClientID:     cfg.GoogleClientId,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       googleAuthScopes,
		Endpoint:     google.Endpoint,
	}
}

func (gp *googleProvider) InitOAuth(role model.UserRole, cbUrl string, force bool) (string, error) {
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOnline,
	}
	if force {
		opts = append(opts, oauth2.ApprovalForce)
	}

	config := gp.config()
	//config.RedirectURL = cbUrl

	u := config.AuthCodeURL(string(role), opts...)
	//u := config.AuthCodeURL("", opts...)
	au, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	//query := au.Query()
	//query.Set("role", string(role))
	//au.RawQuery = query.Encode()

	fmt.Println(au.String())

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
