package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"taxi-service/controllers"
	"taxi-service/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMonitorarCorridaHandler(t *testing.T) {
	app := fiber.New()
	controller := controllers.NewCorridaController()
	app.Post("/corrida/monitorar", controller.MonitorarCorrida)

	body := `{"MotoristaID":1,"PassageiroID":2,"TempoEstimado":20,"TempoDecorrido":25,"Preco":100.0,"Status":"em_andamento"}`
	req := httptest.NewRequest("POST", "/corrida/monitorar", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)

	var corrida models.Corrida
	json.NewDecoder(resp.Body).Decode(&corrida)
	assert.Equal(t, models.StatusAtrasado, corrida.Status)
}

func TestFinalizarCorridaHandler_Antecedencia(t *testing.T) {
	app := fiber.New()
	controller := controllers.NewCorridaController()
	app.Post("/corrida/finalizar", controller.FinalizarCorrida)

	body := `{"MotoristaID":1,"PassageiroID":2,"TempoEstimado":20,"TempoDecorrido":15,"Preco":100.0,"Status":"em_andamento"}`
	req := httptest.NewRequest("POST", "/corrida/finalizar", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var corrida models.Corrida
	json.NewDecoder(resp.Body).Decode(&corrida)
	assert.Equal(t, models.StatusConcluidaAntecedencia, corrida.Status)
	assert.Equal(t, 110.0, corrida.Preco)
	assert.True(t, corrida.BonusAplicado)
}

func TestFinalizarCorridaHandler_NoTempo(t *testing.T) {
	app := fiber.New()
	controller := controllers.NewCorridaController()
	app.Post("/corrida/finalizar", controller.FinalizarCorrida)

	body := `{"MotoristaID":1,"PassageiroID":2,"TempoEstimado":20,"TempoDecorrido":20,"Preco":100.0,"Status":"em_andamento"}`
	req := httptest.NewRequest("POST", "/corrida/finalizar", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var corrida models.Corrida
	json.NewDecoder(resp.Body).Decode(&corrida)
	assert.Equal(t, models.StatusConcluidaNoTempo, corrida.Status)
	assert.Equal(t, 100.0, corrida.Preco)
	assert.False(t, corrida.BonusAplicado)
}



func TestAvaliarCorridaHandler(t *testing.T) {
	app := fiber.New()
	controller := controllers.NewCorridaController()
	app.Post("/corridas/:id/avaliar", controller.AvaliarCorrida)

	// Primeiro: adicionar uma corrida à slice global
	// (poderia ser via POST /corridas, mas aqui é direto)
	model := models.Corrida{
		ID:          99,
		MotoristaID: 1,
		Preco:       50.0,
		Status:      "em_andamento",
	}
	// Adiciona na slice (como mock)
	// Se você tiver um método oficial para isso, use ele
	controllersCorridaService := controller
	controllersCorridaService.Service().AdicionarCorrida(model)

	// Enviar a avaliação
	body := `{"nota": 4}`
	req := httptest.NewRequest("POST", "/corridas/99/avaliar", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}