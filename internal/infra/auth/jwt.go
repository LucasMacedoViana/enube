package auth

import (
	"enube/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateJWT(name string) (string, error) {
	claims := CustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.JWTSecret))
}
