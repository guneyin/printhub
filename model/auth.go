package model

import "encoding/json"

type AuthLoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewAuthLoginRequest(d []byte) (*AuthLoginRequest, error) {
	r := &AuthLoginRequest{}
	err := json.Unmarshal(d, r)
	return r, err
}
