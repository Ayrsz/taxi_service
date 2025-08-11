package controllers

import (
	"taxi-service/models"
	"taxi-service/services"

	"github.com/gofiber/fiber/v2"
	"strconv"
)

// CorridaController gerencia as requisições HTTP para corridas.
type CorridaController struct {
	service *services.CorridaService
}

// NewCorridaController cria uma nova instância de CorridaController.
func NewCorridaController(service *services.CorridaService) *CorridaController {
	return &CorridaController{service: service}
}

// CriarCorrida (POST /corrida) cria uma nova corrida.
func (cc *CorridaController) CriarCorrida(c *fiber.Ctx) error {
	var corridaInput models.Corrida
	if err := c.BodyParser(&corridaInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Não foi possível decodificar o corpo da requisição"})
	}

	if corridaInput.PassageiroID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "PassageiroID é obrigatório"})
	}

	corrida, err := cc.service.CriarNovaCorrida(corridaInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	cc.service.AdicionarCorrida(corridaInput)

	return c.Status(fiber.StatusCreated).JSON(corrida)
}

// GetCorrida (GET /corrida/:id) busca o status de uma corrida.
func (cc *CorridaController) GetCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	corrida, err := cc.service.GetCorridaPorID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(corrida)
}


// MonitorarCorrida (POST /corrida/monitorar) monitora uma corrida.
func (cc *CorridaController) MonitorarCorrida(c *fiber.Ctx) error {
	// A lógica de monitoramento agora será feita pelo frontend buscando o status da corrida.
	return c.SendStatus(fiber.StatusOK)
}

// AceitarCorrida (PUT /corrida/:id/aceitar) permite que um motorista aceite uma corrida.
func (cc *CorridaController) AceitarCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	var body struct {
		MotoristaID int `json:"motoristaId"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Corpo da requisição inválido"})
	}

	if err := cc.service.AceitarCorrida(id, body.MotoristaID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// AtualizarPosicao (PUT /corrida/:id/posicao) atualiza a posição do motorista.
func (cc *CorridaController) AtualizarPosicao(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	var body struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Corpo da requisição inválido"})
	}

	if err := cc.service.AtualizarPosicao(id, body.Lat, body.Lng); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// CancelarCorrida (POST /corrida/:id/cancelar) cancela uma corrida.
func (cc *CorridaController) CancelarCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	if err := cc.service.CancelarCorrida(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// FinalizarCorrida (POST /corrida/:id/finalizar) finaliza uma corrida.
func (cc *CorridaController) FinalizarCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID da corrida inválido"})
	}

	if err := cc.service.FinalizarCorrida(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (cc *CorridaController) AvaliarCorrida(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var input struct {
		Nota int `json:"nota"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	if err := services.AvaliarCorrida(id, input.Nota); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "avaliado"})
}

func (cc *CorridaController) ListarCorridas(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(services.GetCorridas())
}

func (cc *CorridaController) Service() *services.CorridaService {
	return cc.service
}