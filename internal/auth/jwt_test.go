package auth

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuth(t *testing.T) {
	auth := TestAuthenticator{}

	t.Run("return a generated token after passing actual claims", func(t *testing.T) {

		token, err := auth.GenerateToken(testClaims)

		if err != nil {
			t.Errorf("error gottent from generate token %s", err.Error())
		}

		if token == "" {
			t.Error("token string is empty")
		}
	})

	t.Run("return an error after passing nil as claims", func(t *testing.T) {
		token, err := auth.GenerateToken(nil)

		if err == nil {
			t.Fatalf("expecting an actual error %s", err)
		}

		if token != "" {
			t.Fatalf("expecting an empty string. but found %s", token)
		}
	})

	t.Run("validate token generated is correct", func(t *testing.T) {

		token, err := auth.GenerateToken(testClaims)

		jwtToken, err := auth.ValidateToken(token)

		if err != nil {
			t.Fatalf("unexpected error return after validating token %s", err)
		}

		if jwtToken == nil {
			t.Fatal("unexpected nil return for jwtToken")
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		id, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

		if err != nil {
			t.Fatalf("unexpected error while parsing int %s", err)
		}

		if id != 42 {
			t.Fatalf("token generated is wrong. expected %v but found %v", 42, id)
		}
	})

	t.Run("validate token generated is wrong", func(t *testing.T) {

		token, err := auth.GenerateToken(nil)

		jwtToken, err := auth.ValidateToken(token)

		if err == nil {
			t.Fatalf("expected an error %s but err is nil", err)
		}

		if jwtToken != nil {
			t.Fatal("jwtToken should be nil")
		}

	})
}
