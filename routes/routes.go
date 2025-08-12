package routes

import (
	handlers "github.com/gjcms/taxi_service/controllers" // Importe seus handlers
	"github.com/gjcms/taxi_service/middlewares"         // Importe seu middleware
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas as rotas da aplicação Fiber
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1") // Grupo de rotas com prefixo /api/v1

	// Rotas de Autenticação (não precisam de JWT)
	api.Post("/register", handlers.RegisterUser)
	api.Post("/login", handlers.Login)

	// Grupo de rotas protegidas por autenticação JWT
	protected := api.Group("/", middlewares.AuthRequired()) // <-- CORRIGIDO AQUI

	// Rotas para Motoristas (protegidas)
	driver := protected.Group("/drivers/:driverID")
	// Rota para atualização de localização e verificação de chegada (Cenário 1)
	driver.Post("/location", handlers.UpdateDriverLocation)
	// Rota para obter histórico de corridas (Cenários 3 e 4)
	driver.Get("/rides/history", handlers.GetDriverRideHistory)

	// Rotas para Corridas (protegidas)
	rides := protected.Group("/rides")
	rides.Post("/", handlers.RequestRide) // Exemplo: passageiro solicita corrida

	// Rota para motorista aceitar corrida (Cenário 5)
	rides.Post("/:id/accept", handlers.AcceptRide)
	// Rota para motorista finalizar corrida (Cenário 3)
	rides.Post("/:id/complete", handlers.CompleteRide)
	// Rota para motorista cancelar corrida (Cenário 2 - primeira etapa de confirmação)
	rides.Post("/:id/cancel", handlers.CancelRide)
	// Você pode ter uma rota de confirmação de cancelamento separada se a lógica for complexa
	// rides.Post("/:id/confirm_cancel", handlers.ConfirmCancelRide)
}