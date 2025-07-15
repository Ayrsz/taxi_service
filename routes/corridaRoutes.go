// routes/corrida_routes.go
package routes

import (
	"taxi_service/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupCorridaRoutes(api fiber.Router) {
	corridaController := controllers.NewCorridaController()

	corridaRoutes := api.Group("/corridas")

	// Rota para cancelar uma corrida
	// O método é Patch e o caminho é "/:id/cancelar"
	corridaRoutes.Patch("/:id/cancelar", corridaController.CancelarCorrida)
}