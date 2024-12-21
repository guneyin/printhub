package auth

import (
	"github.com/guneyin/printhub/model"
)

type OAuthUser struct {
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Gender    string `json:"gender,omitempty"`
	UserID    string `json:"userID,omitempty"`
	AvatarURL string `json:"avatarURL,omitempty"`
	Location  string `json:"location,omitempty"`
	Link      string `json:"link,omitempty"`
}

func (o *OAuthUser) ToUser() *model.User {
	return &model.User{
		Email:     o.Email,
		Name:      o.Email,
		Password:  "",
		AvatarURL: o.AvatarURL,
	}
}
