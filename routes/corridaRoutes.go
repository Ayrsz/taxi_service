package routes

import (
	"taxi_service/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupCorridaRoutes configura as rotas relacionadas a corridas.
func SetupCorridaRoutes(api fiber.Router) {
	corridaController := controllers.NewCorridaController()

	corridaRoutes := api.Group("/corridas")

	corridaRoutes.Patch("/:id/cancelar", corridaController.CancelarCorrida)

	corridaRoutes.Post("/verificar-tempo", corridaController.VerificarTempo)
}
