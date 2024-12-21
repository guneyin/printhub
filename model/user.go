package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	u.Password = u.PasswordHash()

	if u.Role == "" {
		u.Role = UserRoleClient
	}

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Password = u.PasswordHash()
	return nil
}

func (u *User) Safe() *User {
	user := *u
	user.Password = ""
	return &user
}

func (u *User) JSON() []byte {
	user := *u
	user.Password = ""
	b, _ := json.Marshal(&user)
	return b
}

func (u *User) PasswordHash() string {
	if u.Password == "" {
		return ""
	}

	hashed, _ := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	return string(hashed)
}

func (u *User) IsActivated() bool {
	return u.Active
}

func (u *User) IsValid(role UserRole) bool {
	if usr := u; usr != nil {
		return usr.Role == role
	}
	return false
}

func NewUserRole(s string) (UserRole, error) {
	switch UserRole(s) {
	case UserRoleAdmin, UserRoleTenant, UserRoleClient:
		return UserRole(s), nil
	default:
		return "", fmt.Errorf("invalid user role: %s", s)
	}
}
