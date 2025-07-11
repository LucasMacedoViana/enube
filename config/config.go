package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config estrutura com todas as variáveis carregadas do .env
type Config struct {
	JWTSecret     string
	DatabaseDSN   string
	AdminName     string
	AdminPassword string
	Port          string
}

// Env é a instância global da configuração
var Env *Config

// LoadConfig carrega as variáveis de ambiente do .env e injeta na struct
func LoadConfig() {
	// Carrega o .env se estiver presente (ignora erro se já estiver setado no ambiente)
	_ = godotenv.Load()

	Env = &Config{
		JWTSecret:     os.Getenv("JWT_SECRET"),
		DatabaseDSN:   os.Getenv("DATABASE_DSN"),
		AdminName:     os.Getenv("ADMIN_NAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		Port:          getOrDefault("PORT", "3000"),
	}

	// Valida obrigatoriedade de variáveis essenciais
	if Env.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET não definido no .env")
	}
	if Env.DatabaseDSN == "" {
		log.Fatal("❌ DATABASE_DSN não definido no .env")
	}
}

// getOrDefault busca a env var ou retorna valor padrão
func getOrDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
