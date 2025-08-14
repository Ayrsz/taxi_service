package controllers

import (
	"taxi-service/models"
	"taxi-service/services"

	"github.com/gofiber/fiber/v2"
	"strconv"
)


const (
	// Status originais (mantidos para compatibilidade)
	NoBodyRequestError  = "Não foi possível decodificar o corpo da requisição"
	NoIDRequestedError  = "PassageiroID é obrigatório"
	InvalidIDError = "ID da corrida inválido"
	InvalidRiderIDError = "ID do motorista inválido"

	FinishedRide = "Corrida finalizada"
	StartedRide = "Corrida iniciada"
	CancelRide = "Corrida cancelada"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": NoBodyRequestError})
	}

	if corridaInput.PassageiroID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": NoIDRequestedError})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": InvalidIDError})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": InvalidIDError})
	}

	var body struct {
		MotoristaID int `json:"motoristaId"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": NoBodyRequestError})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": InvalidIDError})
	}

	var body struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": NoBodyRequestError})
	}

	if err := cc.service.AtualizarPosicao(id, body.Lat, body.Lng); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}


// FinalizarCorrida (POST /corrida/:id/finalizar) finaliza uma corrida.
func (cc *CorridaController) FinalizarCorrida(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": InvalidIDError})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": InvalidIDError})
	}

	var input struct {
		Nota int `json:"nota"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": NoBodyRequestError})
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

type cancelamentoMotoristaRequest struct {
	MotoristaID string `json:"motorista_id"`
}

// CancelarCorridaPeloMotorista lida com a requisição de cancelamento de uma corrida por um motorista.
func (cc *CorridaController) CancelarCorridaPeloMotorista(c *fiber.Ctx) error {
    corridaID, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": InvalidIDError,
        })
    }

    var req cancelamentoMotoristaRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": NoBodyRequestError,
        })
    }

    if req.MotoristaID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": InvalidRiderIDError,
        })
    }

    err = cc.service.CancelarCorridaPeloMotorista(corridaID, req.MotoristaID)
    if err != nil {
        // --- MUDANÇA PRINCIPAL AQUI ---
        // Erros de regra de negócio (como "não pode cancelar corrida finalizada")
        // devem retornar um status 400, não 500.
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": CancelRIde,
    })
}

func (cc *CorridaController) IniciarCorrida(c *fiber.Ctx) error {
    corridaID, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": InvalidIDError,
        })
    }

    err = cc.service.IniciarCorrida(corridaID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": StartedRide,
    })
}
