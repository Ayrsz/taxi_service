package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Biblioteca para carregar .env, instale com "go get github.com/joho/godotenv"
)

// AppConfig pode armazenar as configurações da sua aplicação
type AppConfig struct {
	DatabaseURL string
	AppPort     string
	JWTSecret   string
}

// GlobalConfig é uma instância global para acessar as configurações
var GlobalConfig AppConfig

// LoadConfig carrega as variáveis de ambiente
func LoadConfig() {
	// Carrega o arquivo .env, se existir. Para produção, geralmente usa variáveis de ambiente diretas.
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using environment variables directly.")
		}
	}

	GlobalConfig.DatabaseURL = os.Getenv("DATABASE_URL")
	if GlobalConfig.DatabaseURL == "" {
		// Use um valor padrão ou retorne um erro se a URL do DB for crítica
		log.Println("DATABASE_URL not set, using default for development: file::memory:?cache=shared")
		GlobalConfig.DatabaseURL = "file::memory:?cache=shared" // SQLite em memória para dev rápido ou testes
	}

	GlobalConfig.AppPort = os.Getenv("APP_PORT")
	if GlobalConfig.AppPort == "" {
		GlobalConfig.AppPort = "3000" // Porta padrão
	}

	GlobalConfig.JWTSecret = os.Getenv("JWT_SECRET")
	if GlobalConfig.JWTSecret == "" {
		log.Println("JWT_SECRET not set. Using a default secret for development. NOT FOR PRODUCTION!")
		GlobalConfig.JWTSecret = "supersecretjwtkeyforgoapp" // Secret de exemplo, MUDE PARA PRODUÇÃO!
	}

	log.Println("Configuration loaded successfully.")
}