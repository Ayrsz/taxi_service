package config

import (
	"os"
	"testing"
)

// SetupTestEnv sets up the test environment
func SetupTestEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "test.db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DRIVER", "sqlite")
}

// CleanupTestEnv cleans up after tests
func CleanupTestEnv(t *testing.T) {
	// Remove test database file
	if err := os.Remove("test.db"); err != nil {
		t.Logf("Failed to remove test database: %v", err)
	}
}