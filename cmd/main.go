package main

import (
	"enube/config"
	importador "enube/internal/infra/import"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"log"

	_ "enube/docs"            // Importa documentação Swagger
	"enube/internal/handler"  // Rotas da aplicação
	"enube/internal/infra/db" // Conexão e migração do banco
)

// @title       Enube API
// @version     1.0
// @description API para importação, dashboard e autenticação
// @host        localhost:3000
// @BasePath    /
// @schemes     http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Conecta ao banco e aplica migrations
	db.ConnectAndMigrate()

	// Inicializa o app Fiber
	app := fiber.New()

	// Endpoint da documentação Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Define as rotas principais
	handler.SetupRoutes(app)
	//
	// importador
	filePath := "data/excel1.xlsx"
	if err := importador.ImportFromExcel(filePath); err != nil {
		log.Fatalf("Erro ao importar arquivo: %v", err)
	}

	// Inicia o servidor
	port := config.Env.Port
	log.Printf("🚀 Servidor rodando em http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

func init() {
	// Carrega variáveis de ambiente
	config.LoadConfig()
}
