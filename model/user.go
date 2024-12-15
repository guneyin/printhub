package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleTenant UserRole = "tenant"
	UserRoleClient UserRole = "client"
)

type User struct {
	gorm.Model `json:"-"`
	UUID       string   `json:"uuid" gorm:"index"`
	Role       UserRole `json:"role" gorm:"index:,unique,composite:idx_user_unique"`
	Email      string   `json:"email" gorm:"index:,unique,composite:idx_user_unique"`
	Name       string   `json:"name"`
	Password   string   `json:"password,omitempty"`
	AvatarURL  string   `json:"avatarURL"`
	Active     bool     `json:"active"`
}

type UserList []User

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New().String()

	if u.Role == "" {
		u.Role = UserRoleClient
	}

	return nil
}

func (u *User) JSON() []byte {
	user := *u
	user.Password = ""
	b, _ := json.Marshal(&user)
	return b
}

func (u *User) IsActivated() bool {
	return u.Active
}

func NewUserRole(s any) (UserRole, error) {
	if role, ok := s.(UserRole); ok {
		return role, nil
	}
	return "", fmt.Errorf("invalid user role: %v", s)
}
