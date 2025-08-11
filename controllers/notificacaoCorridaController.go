package controllers

import (
	"strconv"
	"strings"
	"taxi-service/models"
	"taxi-service/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ListNotificacoesCorrida - Lista todas as notificações
func ListNotificacoesCorrida(c *fiber.Ctx) error {
	notificacoes, err := services.ListNotificacoesCorrida()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notificacoes",
		})
	}
	return c.Status(fiber.StatusOK).JSON(notificacoes)
}

// GetNotificacaoCorrida - Busca notificação por ID
func GetNotificacaoCorrida(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	notificacaoID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	notificacao, err := services.GetNotificacaoCorrida(uint(notificacaoID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notificacao not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notificacao",
		})
	}

	return c.Status(fiber.StatusOK).JSON(notificacao)
}

// CreateNotificacaoCorrida - Cria nova notificação para motorista
func CreateNotificacaoCorrida(c *fiber.Ctx) error {
	notificacao := new(models.NotificacaoCorrida)
	if err := c.BodyParser(notificacao); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validar dados obrigatórios
	if notificacao.MotoristaID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "MotoristaID is required",
		})
	}
	if notificacao.CorridaID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "CorridaID is required",
		})
	}
	if notificacao.PassageiroNome == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "PassageiroNome is required",
		})
	}
	if notificacao.Valor <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Valor must be greater than 0",
		})
	}

	err := services.CreateNotificacaoCorrida(notificacao)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create notificacao",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(notificacao)
}

// GetNotificacoesPendentesParaMotorista - Busca notificações pendentes para um motorista específico
func GetNotificacoesPendentesParaMotorista(c *fiber.Ctx) error {
	motoristaIDParam := c.Params("motoristaID")

	if motoristaIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "MotoristaID is required",
		})
	}

	motoristaID, err := strconv.ParseUint(motoristaIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid MotoristaID format",
		})
	}

	notificacoes, err := services.GetNotificacoesPendentesParaMotorista(uint(motoristaID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pending notificacoes",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"motorista_id":         motoristaID,
		"pending_count":        len(notificacoes),
		"pending_notificacoes": notificacoes,
	})
}

// AceitarNotificacaoCorrida - Aceita uma notificação de corrida
func AceitarNotificacaoCorrida(c *fiber.Ctx) error {
	notificacaoID := c.Params("id")
	motoristaIDParam := c.Params("motoristaID")

	if notificacaoID == "" || motoristaIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "NotificacaoID and MotoristaID are required",
		})
	}

	nID, err := strconv.ParseUint(notificacaoID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid NotificacaoID format",
		})
	}

	mID, err := strconv.ParseUint(motoristaIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid MotoristaID format",
		})
	}

	err = services.AceitarNotificacaoCorrida(uint(nID), uint(mID))
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return c.Status(fiber.StatusGone).JSON(fiber.Map{
				"error": "Notificacao expired",
			})
		}
		if strings.Contains(err.Error(), "already processed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Notificacao already processed",
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notificacao not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to accept notificacao",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Notificacao accepted successfully",
		"notificacao_id": nID,
		"motorista_id":   mID,
	})
}

// RecusarNotificacaoCorrida - Recusa uma notificação de corrida
func RecusarNotificacaoCorrida(c *fiber.Ctx) error {
	notificacaoID := c.Params("id")
	motoristaIDParam := c.Params("motoristaID")

	if notificacaoID == "" || motoristaIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "NotificacaoID and MotoristaID are required",
		})
	}

	nID, err := strconv.ParseUint(notificacaoID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid NotificacaoID format",
		})
	}

	mID, err := strconv.ParseUint(motoristaIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid MotoristaID format",
		})
	}

	err = services.RecusarNotificacaoCorrida(uint(nID), uint(mID))
	if err != nil {
		if strings.Contains(err.Error(), "already processed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Notificacao already processed",
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notificacao not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refuse notificacao",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Notificacao refused successfully",
		"notificacao_id": nID,
		"motorista_id":   mID,
	})
}

// ExpirarNotificacoesVencidas - Marca como expiradas as notificações que passaram do tempo limite
func ExpirarNotificacoesVencidas(c *fiber.Ctx) error {
	err := services.ExpirarNotificacoesVencidas()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to expire notificacoes",
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Expired notificacoes processed successfully",
		"processed_at": time.Now(),
	})
}

// GetHistoricoNotificacoesMotorista - Busca histórico de notificações de um motorista
func GetHistoricoNotificacoesMotorista(c *fiber.Ctx) error {
	motoristaIDParam := c.Params("motoristaID")

	if motoristaIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "MotoristaID is required",
		})
	}

	motoristaID, err := strconv.ParseUint(motoristaIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid MotoristaID format",
		})
	}

	historico, err := services.GetHistoricoNotificacoesMotorista(uint(motoristaID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch historico",
		})
	}

	// Organizar histórico por status
	aceitas := 0
	recusadas := 0
	expiradas := 0
	pendentes := 0

	for _, notif := range historico {
		switch notif.Status {
		case models.NotificacaoAceita:
			aceitas++
		case models.NotificacaoRecusada:
			recusadas++
		case models.NotificacaoExpirada:
			expiradas++
		case models.NotificacaoPendente:
			pendentes++
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"motorista_id":    motoristaID,
		"total_count":     len(historico),
		"aceitas_count":   aceitas,
		"recusadas_count": recusadas,
		"expiradas_count": expiradas,
		"pendentes_count": pendentes,
		"historico":       historico,
	})
}

// DeleteNotificacaoCorrida - Remove uma notificação (para limpeza de dados antigos)
func DeleteNotificacaoCorrida(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	notificacaoID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	err = services.DeleteNotificacaoCorrida(uint(notificacaoID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notificacao not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notificacao",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notificacao deleted successfully",
		"id":      notificacaoID,
	})
}

// UpdateNotificacaoStatus - Atualiza status de uma notificação (função auxiliar)
func UpdateNotificacaoStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	newStatus := c.Query("status")

	if id == "" || newStatus == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID and status are required",
		})
	}

	notificacaoID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	// Validar status
	validStatuses := []string{
		"pendente",
		"aceita",
		"recusada",
		"expirada",
	}

	statusValido := false
	for _, validStatus := range validStatuses {
		if validStatus == newStatus {
			statusValido = true
			break
		}
	}

	if !statusValido {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Valid values: pendente, aceita, recusada, expirada",
		})
	}

	// Buscar a notificação atual
	notificacao, err := services.GetNotificacaoCorrida(uint(notificacaoID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notificacao not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notificacao",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Status update endpoint - implementation needed",
		"current":    notificacao,
		"new_status": newStatus,
		"suggestion": "Use specific endpoints: /accept, /refuse, or /expire",
	})
}
