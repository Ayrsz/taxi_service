package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"taxi-service/controllers"
	"taxi-service/models"
	"taxi-service/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// setupApp é uma função helper para inicializar o app, serviço e controller para cada teste.
// Isso garante que os testes sejam independentes um do outro.
func setupApp() (*fiber.App, *services.CorridaService) {
	app := fiber.New()
	corridaService := services.NewCorridaService()
	controller := controllers.NewCorridaController(corridaService)

	// Configuração de todas as rotas necessárias para o teste de cancelamento
	app.Post("/corrida", controller.CriarCorrida)
	app.Post("/corrida/:id/aceitar", controller.AceitarCorrida)
	app.Post("/corrida/:id/finalizar", controller.FinalizarCorrida)
	app.Post("/corrida/:id/cancelar/motorista", controller.CancelarCorridaPeloMotorista)

	return app, corridaService
}

// Teste principal para o fluxo de cancelamento de corrida.
func TestCancelamentoDeCorrida(t *testing.T) {

	// Cenário: Motorista cancela uma corrida pendente (aceita, mas não iniciada)
	t.Run("Motorista pode cancelar uma corrida que acabou de aceitar", func(t *testing.T) {
		// --- ARRANGE (Preparação) ---
		app, corridaService := setupApp()
		corrida, _ := corridaService.CriarNovaCorrida(models.Corrida{PassageiroID: 1})
		motoristaID := 123
		corridaService.AceitarCorrida(corrida.ID, motoristaID) // Status agora é 'motorista_encontrado'

		// --- ACT (Ação) ---
		payload := map[string]string{"motorista_id": fmt.Sprint(motoristaID)}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/corrida/%d/cancelar/motorista", corrida.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// --- ASSERT (Verificação) ---
		assert.Equal(t, http.StatusOK, resp.StatusCode, "A resposta da API deve ser 200 OK")
		corridaCancelada, _ := corridaService.GetCorridaPorID(corrida.ID)
		assert.Equal(t, models.StatusCanceladaPeloMotorista, corridaCancelada.Status, "O status da corrida deve ser 'cancelada pelo motorista'")
	})

	// Cenário: Tentativa de cancelamento de corrida em andamento pelo motorista
	t.Run("Motorista não pode cancelar uma corrida em andamento", func(t *testing.T) {
		// --- ARRANGE ---
		app, corridaService := setupApp()
		corrida, _ := corridaService.CriarNovaCorrida(models.Corrida{PassageiroID: 1})
		motoristaID := 456
		corridaService.AceitarCorrida(corrida.ID, motoristaID)
		// Forçamos o status para "em_andamento" para simular o cenário
		corrida.Status = models.StatusEmAndamento

		// --- ACT ---
		payload := map[string]string{"motorista_id": fmt.Sprint(motoristaID)}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/corrida/%d/cancelar/motorista", corrida.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// --- ASSERT ---
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "A API deve retornar um erro para corrida em andamento")
		corridaAposTentativa, _ := corridaService.GetCorridaPorID(corrida.ID)
		assert.Equal(t, models.StatusEmAndamento, corridaAposTentativa.Status, "O status da corrida não deve mudar")
	})

	// Cenário: Tentativa de cancelamento de corrida já finalizada pelo motorista
	t.Run("Motorista não pode cancelar uma corrida finalizada", func(t *testing.T) {
		// --- ARRANGE ---
		app, corridaService := setupApp()
		corrida, _ := corridaService.CriarNovaCorrida(models.Corrida{PassageiroID: 1})
		motoristaID := 789
		corridaService.AceitarCorrida(corrida.ID, motoristaID)
		corridaService.FinalizarCorrida(corrida.ID) // Finaliza a corrida

		// --- ACT ---
		payload := map[string]string{"motorista_id": fmt.Sprint(motoristaID)}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/corrida/%d/cancelar/motorista", corrida.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// --- ASSERT ---
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "A API deve retornar um erro para corrida finalizada")
		corridaAposTentativa, _ := corridaService.GetCorridaPorID(corrida.ID)
		// O status deve ser um dos de conclusão, e não "cancelada"
		assert.Contains(t, []string{models.StatusConcluidaNoTempo, models.StatusConcluidaAntecedencia, models.StatusAtrasado}, corridaAposTentativa.Status)
	})

	// Este teste é uma validação de segurança e não está no feature, mas é importante.
	t.Run("Motorista não pode cancelar uma corrida de outro motorista", func(t *testing.T) {
		// --- ARRANGE ---
		app, corridaService := setupApp()
		corrida, _ := corridaService.CriarNovaCorrida(models.Corrida{PassageiroID: 1})
		corridaService.AceitarCorrida(corrida.ID, 123) // Corrida do motorista 123

		// --- ACT ---
		payload := map[string]string{"motorista_id": "999"} // Motorista 999 (impostor) tenta cancelar
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/corrida/%d/cancelar/motorista", corrida.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// --- ASSERT ---
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "A API deve retornar um erro quando um motorista não autorizado tenta cancelar")
		corridaAposTentativa, _ := corridaService.GetCorridaPorID(corrida.ID)
		assert.Equal(t, models.StatusMotoristaEncontrado, corridaAposTentativa.Status, "O status da corrida não deve mudar")
	})
}
