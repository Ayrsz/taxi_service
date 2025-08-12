package routes

import (
	"taxi-service/controllers"
	"taxi-service/services"
	"github.com/gofiber/fiber/v2"
)

// SetupCorridaRoutes configura as rotas relacionadas a corridas.
func SetupCorridaRoutes(api fiber.Router, corridaService *services.CorridaService) {
	corridaController := controllers.NewCorridaController(corridaService)

	corridaGroup := api.Group("/corrida")
	corridaGroup.Post("/", corridaController.CriarCorrida)
	corridaGroup.Get("/:id", corridaController.GetCorrida)
	corridaGroup.Post("/monitorar", corridaController.MonitorarCorrida)
	corridaGroup.Put("/:id/aceitar", corridaController.AceitarCorrida)
	corridaGroup.Put("/:id/posicao", corridaController.AtualizarPosicao)
	corridaGroup.Post("/:id/cancelar", corridaController.CancelarCorrida)
	corridaGroup.Post("/:id/finalizar", corridaController.FinalizarCorrida)

	api.Post("/corridas/:id/avaliar", corridaController.AvaliarCorrida)
	api.Post("/corridas", corridaController.CriarCorrida)
	api.Get("/corridas", corridaController.ListarCorridas)

	// Manter a rota OPTIONS para o CORS
	corridaGroup.Options("/monitorar", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
}
