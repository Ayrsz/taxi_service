package controllers

import (
	"strconv"
	"strings"
	"taxi_service/database"
	"taxi_service/models"
	"taxi_service/services"

	"github.com/gofiber/fiber/v2"
)

type CorridaController struct {
	service *services.CorridaService
}

func NewCorridaController() *CorridaController {
	repo := database.NewJSONCorridaRepository("./data/corridas.json")
	service := services.NewCorridaService(repo)
	return &CorridaController{service: service}
}

// handler para a rota de cancelamento.
func (ctrl *CorridaController) CancelarCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	corrida, err := ctrl.service.CancelarCorrida(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(corrida)
}

// handler para a rota de cancelamento.
func (ctrl *CorridaController) CancelarCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	corrida, err := ctrl.service.CancelarCorrida(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(corrida)
}

// POST /corrida/cancelar-por-excesso-tempo
func (cc *CorridaController) CancelarPorExcessoTempo(c *fiber.Ctx) error {
	var corrida models.Corrida
	if err := c.BodyParser(&corrida); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if corrida.TempoDecorrido-corrida.TempoEstimado > 15 {
		corrida.Status = models.StatusCanceladaPorExcessoTempo
		services.NotificarMotorista(c, corrida.MotoristaID, "Corrida cancelada por excesso de tempo")
	}
	return c.Status(fiber.StatusOK).JSON(corrida)
}

// Handler para a rota de verificação de tempo do sistema.
func (ctrl *CorridaController) VerificarTempo(c *fiber.Ctx) error {
	if err := ctrl.service.VerificarCorridasPendentesPorTimeout(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Falha ao executar a verificação de tempo",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Verificação de tempo de corridas pendentes concluída.",
	})

}
