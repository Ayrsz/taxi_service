package test

import (
	"taxi-service/models"
	"taxi-service/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAplicarBonus(t *testing.T) {
	corrida := &models.Corrida{Preco: 100.0}
	services.AplicarBonusSTUB(corrida)
	assert.Equal(t, 110.0, corrida.Preco)
}

func TestVerificarTempoCorrida_Atrasado(t *testing.T) {
	service := services.CorridaServiceSTUB{}
	corrida := models.Corrida{
		TempoEstimado:  20,
		TempoDecorrido: 25,
		Status:         models.StatusEmAndamento,
	}

	service.VerificarTempoCorridaSTUB(nil, &corrida)
	assert.Equal(t, models.StatusAtrasado, corrida.Status)
}

func TestVerificarTempoCorrida_Cancelada(t *testing.T) {
	service := services.CorridaServiceSTUB{}
	corrida := models.Corrida{
		TempoEstimado:  20,
		TempoDecorrido: 36,
		Status:         models.StatusEmAndamento,
	}

	service.VerificarTempoCorridaSTUB(nil, &corrida)
	assert.Equal(t, models.StatusCanceladaPorExcessoTempo, corrida.Status)
	assert.Equal(t, models.StatusCanceladaPorExcessoTempo, corrida.Status)
}

func TestFinalizarCorrida_Antecedencia(t *testing.T) {
	service := services.CorridaServiceSTUB{}
	corrida := models.Corrida{
		TempoEstimado:  20,
		TempoDecorrido: 15,
		Preco:          100.0,
		Status:         models.StatusEmAndamento,
	}

	service.FinalizarCorridaSTUB(nil, &corrida)
	assert.Equal(t, models.StatusConcluidaAntecedencia, corrida.Status)
	assert.Equal(t, 110.0, corrida.Preco)
	assert.True(t, corrida.BonusAplicado)
	assert.True(t, corrida.BonusAplicado)
}

func TestFinalizarCorrida_NoTempoSTUB(t *testing.T) {
	service := services.CorridaServiceSTUB{}
	corrida := models.Corrida{
		TempoEstimado:  20,
		TempoDecorrido: 20,
		Preco:          100.0,
		Status:         models.StatusEmAndamento,
	}

	service.FinalizarCorridaSTUB(nil, &corrida)
	assert.Equal(t, models.StatusConcluidaNoTempo, corrida.Status)
	assert.Equal(t, 100.0, corrida.Preco)
	assert.False(t, corrida.BonusAplicado)
	assert.False(t, corrida.BonusAplicado)
}
