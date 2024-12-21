package mw

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"strconv"
	"time"
)

type jwtPayload struct {
	UserID   string         `json:"user_id"`
	UserName string         `json:"user_name"`
	UserRole model.UserRole `json:"user_role"`
}

func genJWT(sess *model.Session) (string, error) {
	cfg := market.Get().Config

	exp, _ := strconv.Atoi(cfg.JWTExp)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":       time.Now().Add(time.Hour * 24 * time.Duration(exp)).Unix(),
		"user_id":   sess.UserId,
		"user_name": sess.UserEmail,
		"user_role": sess.UserRole,
	})

	token, err := t.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func parseJWT(token string) (*jwt.Token, error) {
	cfg := market.Get().Config

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(cfg.JWTSecret), nil
	})
}

func verifyJWT(token string) (*jwtPayload, error) {
	parsed, err := parseJWT(token)
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt claim error: %v", err)
	}

	payload := &jwtPayload{}
	payload.UserID, ok = claims["user_id"].(string)
	payload.UserName, ok = claims["user_name"].(string)
	payload.UserRole, ok = claims["user_role"].(model.UserRole)

	if !ok {
		return nil, errors.New("something went wrong")
	}

	return payload, nil
}
