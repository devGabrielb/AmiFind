package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenDetails struct {
	secret string
}

func NewToken(env string) *tokenDetails {
	return &tokenDetails{secret: env}
}

func (td *tokenDetails) GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(td.secret))
	if err != nil {
		return "", err
	}
	return t, nil
}
