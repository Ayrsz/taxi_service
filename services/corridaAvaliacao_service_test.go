package services

import (
	"testing"
	"taxi-service/models"
)

func TestAvaliarCorrida_Sucesso(t *testing.T) {
	corridas = []models.Corrida{
		{ID: 10, MotoristaID: 999},
	}

	err := AvaliarCorrida(10, 5)
	if err != nil {
		t.Fatalf("Esperava sucesso, mas deu erro: %v", err)
	}

	if corridas[0].Avaliacao == nil || *corridas[0].Avaliacao != 5 {
		t.Errorf("Esperava nota 5, recebeu: %v", corridas[0].Avaliacao)
	}
}

func TestAvaliarCorrida_CorridaNaoEncontrada(t *testing.T) {
	corridas = []models.Corrida{} // vazio

	err := AvaliarCorrida(999, 4)
	if err == nil {
		t.Fatalf("Esperava erro por corrida inexistente, mas foi nil")
	}
}
