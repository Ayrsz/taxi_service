package services

import (
	"fmt"
	"taxi-service/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CorridaServiceSTUB struct {
	Corridas []models.Corrida // Mock de "banco de dados" em memória
}

// NotificacaoService é responsável por enviar notificações
func NotificarPassageiroSTUB(_ *fiber.Ctx, passageiroID int, mensagem string) error {
	fmt.Printf("[Stub] Notificação ao passageiro %d: %s\n", passageiroID, mensagem)
	return nil
}

func NotificarMotoristaSTUB(_ *fiber.Ctx, motoristaID int, mensagem string) error {
	fmt.Printf("[Stub] Notificação ao motorista %d: %s\n", motoristaID, mensagem)
	return nil
}

func AplicarBonusSTUB(corrida *models.Corrida) {
	// Aplica 10% de bônus sobre o preço da corrida, se o campo existir
	if corrida != nil {
		corrida.Preco = corrida.Preco + (corrida.Preco * 0.10)
	}

}

func (s *CorridaServiceSTUB) VerificarTempoCorridaSTUB(_ interface{}, corrida *models.Corrida) {
	diferenca := corrida.TempoDecorrido - corrida.TempoEstimado

	if corrida.Status == models.StatusEmAndamento {
		if diferenca > 0 && diferenca <= 15 {
			corrida.Status = models.StatusAtrasado
		} else if diferenca > 15 {
			corrida.Status = models.StatusCanceladaPorExcessoTempo
		}
	}
}

func (s *CorridaServiceSTUB) FinalizarCorridaSTUB(_ interface{}, corrida *models.Corrida) {
	diferenca := corrida.TempoDecorrido - corrida.TempoEstimado
	now := time.Now()
	corrida.DataFim = &now

	if diferenca < 0 {
		corrida.Status = models.StatusConcluidaAntecedencia
		AplicarBonusSTUB(corrida)
		corrida.BonusAplicado = true
	} else if diferenca == 0 || diferenca > 0 && diferenca <= 15 {
		corrida.Status = models.StatusConcluidaNoTempo
	}
}
