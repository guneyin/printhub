package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model `json:"-"`
	UUID       string `json:"uuid" gorm:"index"`
	Email      string `json:"email" gorm:"uniqueIndex"`
	Name       string `json:"name" gorm:"index"`
	Address    string `json:"address"`
	Logo       string `json:"logo"`
}

type TenantList []Tenant

type TenantUser struct {
	gorm.Model
	TenantID uint `gorm:"index"`
	UserID   uint
}

type TenantUserList struct {
	Tenant Tenant
	Users  []UserList
}

func NewTenant(d []byte) (*Tenant, error) {
	t := &Tenant{}
	err := json.Unmarshal(d, t)
	return t, err
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	t.UUID = uuid.New().String()

	return nil
}
