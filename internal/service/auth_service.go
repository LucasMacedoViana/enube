package service

import (
	"enube/internal/infra/auth"
	"github.com/gofiber/fiber/v2"
)

func Authenticate(name, password string) (string, error) {
	// Simulação de login simples (em produção: validar hash da senha)
	if name == "admin" && password == "admin123" {
		return auth.GenerateJWT(name)
	}
	return "", fiber.ErrUnauthorized
}
