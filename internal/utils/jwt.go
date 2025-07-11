package utils

import (
	"enube/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateJWT cria um token JWT com ID e nome do usuário
func GenerateJWT(name, password string) (string, error) {
	secret := config.Env.JWTSecret
	if secret == "" {
		secret = "secret" // fallback seguro para dev
	}

	claims := jwt.MapClaims{
		"name":     name,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
