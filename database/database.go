package database

import (
	"log"

	"github.com/gjcms/taxi_service/config" // Importe seu pacote de configuração
	"github.com/gjcms/taxi_service/models" // Importe seus modelos para automigração
	"gorm.io/driver/postgres"              // Se for usar PostgreSQL em prod, instale com "go get gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"                // Para testes ou dev leve, instale com "go get gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Para controlar os logs do GORM
)

// DB é a instância global do GORM para interagir com o banco de dados
var DB *gorm.DB

// ConnectDB estabelece a conexão com o banco de dados principal
func ConnectDB() {
	var err error
	dsn := config.GlobalConfig.DatabaseURL // Pega a URL do DB das configurações

	// Escolha o driver do banco de dados com base na DSN ou configuração
	if contains(dsn, "postgres") { // Exemplo para PostgreSQL
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Logar consultas SQL em dev
		})
	} else { // Padrão para SQLite (usado em memória para dev/teste se nada for configurado)
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Logar consultas SQL em dev
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully!")

	// Migrar os modelos automaticamente (cria tabelas se não existirem)
	err = DB.AutoMigrate(&models.User{}, &models.Ride{}, &models.DummyUser{}) // Se DummyUser for um modelo ativo
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	log.Println("Database migration completed!")
}

// contains é uma função auxiliar simples para verificar strings
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
