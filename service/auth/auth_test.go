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

	verified, err := validateToken(hashed)
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
	svc.Recover(ctx, email, role)
}

func TestValidate(t *testing.T) {
	ctx := context.Background()
	svc := newService()
	u, err := svc.Validate(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u.Safe())
}

const token = "xDA9HiSLJ6kqGHqZZhJn3k4BONLUIniAJq1Q4SdZmVWAiFqy4RssHST1CI5tk2oFtTMqmirXDMHuiTyGp93qeCWj5NJflxYKmyeQ"

func TestValidateToken(t *testing.T) {
	uid, err := validateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("uuid:", uid)
}
