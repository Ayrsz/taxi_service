package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"taxi_service/models"
	"taxi_service/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// preparar o ambiente de cada teste.
func setupCorridasTestFile(t *testing.T, corridas []models.Corrida) string {
	dataDir := "./data"
	err := os.MkdirAll(dataDir, 0755)
	assert.NoError(t, err)

	filePath := filepath.Join(dataDir, "corridas.json")
	data, err := json.MarshalIndent(corridas, "", "  ")
	assert.NoError(t, err)

	err = os.WriteFile(filePath, data, 0644)
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(dataDir)
	})

	// Retornamos o caminho para o caso de precisarmos ler o arquivo depois
	return filePath
}

// Testa o Cenário: "Motorista cancela uma corrida pendente"
func TestMotoristaCancelaCorridaPendente(t *testing.T) {
	// Arrange
	app := test.SetupTestApp(t)
	corridasIniciais := []models.Corrida{
		{Id: 101, Status: "pendente", CPFMotorista: new(int)},
	}
	*corridasIniciais[0].CPFMotorista = 12345
	setupCorridasTestFile(t, corridasIniciais)

	// Act
	resp := test.MakeRequest(t, app, "PATCH", "/api/corridas/101/cancelar", nil)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var corridaCancelada models.Corrida
	test.ParseResponseBody(t, resp, &corridaCancelada)
	assert.Equal(t, "cancelada", corridaCancelada.Status)
}

// Testa o Cenário: "Sistema cancela corrida por demora na partida do motorista"
func TestSistemaCancelaCorridaPorDemora(t *testing.T) {
	// Arrange
	app := test.SetupTestApp(t)
	horarioAntigo := time.Now().Add(-15 * time.Minute)
	corridasIniciais := []models.Corrida{
		{Id: 101, Status: "pendente", Horario: horarioAntigo},
	}
	filePath := setupCorridasTestFile(t, corridasIniciais)

	// Act
	resp := test.MakeRequest(t, app, "POST", "/api/corridas/verificar-tempo", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	var corridasDepois []models.Corrida
	err = json.Unmarshal(data, &corridasDepois)
	assert.NoError(t, err)

	assert.Equal(t, "cancelada", corridasDepois[0].Status)
}

// Testa o Cenário: "Tentativa de cancelamento de corrida em andamento pelo motorista"
func TestMotoristaTentaCancelarCorridaEmAndamento(t *testing.T) {
	// Arrange
	app := test.SetupTestApp(t)
	corridasIniciais := []models.Corrida{
		{Id: 102, Status: "andamento"},
	}
	setupCorridasTestFile(t, corridasIniciais)

	// Act
	resp := test.MakeRequest(t, app, "PATCH", "/api/corridas/102/cancelar", nil)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errorResponse map[string]string
	test.ParseResponseBody(t, resp, &errorResponse)
	assert.Contains(t, errorResponse["error"], "não é possível cancelar uma corrida que não está em andamento")
}

func TestMotoristaTentaCancelarCorridaFinalizada(t *testing.T) {
	app := test.SetupTestApp(t)
	corridasIniciais := []models.Corrida{
		{Id: 103, Status: "finalizada"},
	}
	setupCorridasTestFile(t, corridasIniciais)
	resp := test.MakeRequest(t, app, "PATCH", "/api/corridas/103/cancelar", nil)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	var errorResponse map[string]string
	test.ParseResponseBody(t, resp, &errorResponse)
	assert.Contains(t, errorResponse["error"], "não é possível cancelar uma corrida que não está em andamento")
}
