package database

import (
	"log"

	"github.com/gjcms/taxi_service/models" // Importe seus modelos
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectTestDB estabelece a conexão com um banco de dados SQLite em memória para testes
func ConnectTestDB() {
	var err error
	// Conecta a um banco de dados SQLite em memória para testes
	// O parâmetro `_journal_mode=WAL` melhora a concorrência e o desempenho em testes
	DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared&_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Logs silenciosos para testes
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	log.Println("Connected to in-memory test database successfully!")

	// Migrar os modelos automaticamente para o DB de teste
	err = DB.AutoMigrate(&models.User{}, &models.Ride{}, &models.DummyUser{})
	if err != nil {
		log.Fatalf("Failed to auto migrate test database: %v", err)
	}
	log.Println("Test database migration completed!")
}