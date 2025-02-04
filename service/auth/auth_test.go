package auth_test

import (
	"context"
	"log"
	"testing"

	"github.com/guneyin/printhub/service/auth"

	"github.com/google/uuid"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Print("Error loading .env file")
	}
	market.InitMarket()
}

func TestToken(t *testing.T) {
	uid := uuid.New()
	hashed, err := auth.GenerateToken(uid.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("hashed:", hashed)

	verified, err := auth.VerifyToken(hashed)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("verified:", verified)

	if uid.String() != verified {
		t.Fatal("uid does not match")
	}
}

func TestForgotPassword(_ *testing.T) {
	ctx := context.Background()
	email := "guneyin@gmail.com"
	role := model.UserRoleAdmin

	svc := auth.GetService()
	svc.RecoverPassword(ctx, email, role)
}

func TestValidate(t *testing.T) {
	ctx := context.Background()
	svc := auth.GetService()

	token := generateTestToken()
	u, err := svc.ValidateUser(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u.Safe())
}

func generateTestToken() string {
	hashed, _ := auth.GenerateToken(uuid.New().String())
	return hashed
}

func TestValidateToken(t *testing.T) {
	token := generateTestToken()
	uid, err := auth.VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("uuid:", uid)
}
