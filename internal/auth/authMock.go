package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "test"

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-iss",
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),
}

type TestAuthenticator struct{}

func (t TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	if claims == nil {
		return "", errors.New("claims can't be nil")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (t TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
