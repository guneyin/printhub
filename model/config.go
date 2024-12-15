package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model `json:"-"`
	Identifier string `json:"identifier" gorm:"index:,unique,composite:idx_config_unique"`
	Module     string `json:"module" gorm:"index:,unique,composite:idx_config_unique"`
	Key        string `json:"key" gorm:"index:,unique,composite:idx_config_unique"`
	Value      string `json:"value"`
}

type ConfigList []Config

func NewConfigList(d []byte) (*ConfigList, error) {
	dcl := &ConfigList{}
	err := json.Unmarshal(d, dcl)
	return dcl, err
}

func (dcl *ConfigList) JSON() []byte {
	b, _ := json.Marshal(dcl)
	return b
}

func (dc *Config) JSON() []byte {
	b, _ := json.Marshal(dc)
	return b
}
