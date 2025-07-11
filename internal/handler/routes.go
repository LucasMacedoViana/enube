package handler

import (
	"enube/internal/handler/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Rota pública
	api.Post("/login", Login)

	// Rotas protegidas por JWT
	protected := api.Use(middleware.JWTProtected())

	// Usuários
	protected.Post("/users", CreateUser)

	// Dashboard
	protected.Get("/dashboard/summary", GetSummary)
	protected.Get("/dashboard/monthly", GetByMonth)
	protected.Get("/dashboard/by-category", GetByCategory)
	protected.Get("/dashboard/by-client", GetByClient)
	protected.Get("/dashboard/by-resource", GetByResource)
}
