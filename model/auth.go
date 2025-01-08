package model

import "encoding/json"

type AuthUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewAuthUserRequest(d []byte) (*AuthUserRequest, error) {
	r := &AuthUserRequest{}
	err := json.Unmarshal(d, r)
	return r, err
}
