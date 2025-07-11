package db

import (
	"enube/internal/domain/model"
	"enube/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// ConnectAndMigrate realiza a conexão com o banco e executa as migrations
func ConnectAndMigrate() {
	start := time.Now()

	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
		log.Fatal("A variável DATABASE_DSN não está definida")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v", err)
	}
	DB = db
	log.Println("✅ Conexão com banco de dados estabelecida.")

	if err := applyMigrations(); err != nil {
		log.Fatalf("Erro ao aplicar AutoMigrate: %v", err)
	}

	createAdminUser()

	elapsed := time.Since(start)
	log.Printf("⏱️ Migração concluída em %s.\n", elapsed)
}

// applyMigrations executa as migrations das tabelas
func applyMigrations() error {
	log.Println("🔁 Aplicando migrations...")
	return DB.AutoMigrate(
		&model.User{},
		&model.Partner{},
		&model.Customer{},
		&model.Subscription{},
		&model.Meter{},
		&model.Product{},
		&model.BillingItem{},
	)
}

// createAdminUser cria o usuário admin se ele não existir
func createAdminUser() {
	name := os.Getenv("ADMIN_NAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if name == "" || password == "" {
		log.Println("⚠️ ADMIN_NAME ou ADMIN_PASSWORD não definidas, pulando criação do admin.")
		return
	}

	var user model.User
	err := DB.Where("name = ?", name).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		hashedPassword, hashErr := utils.HashPassword(password)
		if hashErr != nil {
			log.Fatalf("Erro ao gerar hash da senha do admin: %v", hashErr)
		}

		admin := model.User{
			Name:     name,
			Password: hashedPassword,
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalf("Erro ao criar usuário admin: %v", err)
		}

		log.Println("✅ Usuário admin criado com sucesso.")
	} else if err != nil {
		log.Fatalf("Erro ao verificar existência do usuário admin: %v", err)
	} else {
		log.Println("ℹ️ Usuário admin já existe.")
	}
}
