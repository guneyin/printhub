package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"github.com/joho/godotenv"
	"log/slog"
	"testing"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		slog.Warn("error loading test .env file")
	}
	market.InitMarket()
}

func TestToken(t *testing.T) {
	uid := uuid.New()
	hashed, err := generateToken(uid.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("hashed:", hashed)

	verified, err := verifyToken(hashed)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("verified:", verified)

	if uid.String() != verified {
		t.Fatal("uid does not match")
	}
}

func TestForgotPassword(t *testing.T) {
	ctx := context.Background()
	email := "guneyin@gmail.com"
	role := model.UserRoleAdmin

	svc := newService()
	svc.RecoverPassword(ctx, email, role)
}

func TestValidate(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	token := generateTestToken()
	u, err := svc.ValidateUser(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u.Safe())
}

func generateTestToken() string {
	hashed, _ := generateToken(uuid.New().String())
	return hashed
}

func TestValidateToken(t *testing.T) {
	token := generateTestToken()
	uid, err := verifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("uuid:", uid)
}
