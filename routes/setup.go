package routes

import (
	"taxi-service/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupAppRoutes inicializa todas as rotas da aplicação.
func SetupRoutes(app *fiber.App) {
	// Middlewares
	app.Use(cors.New())
	app.Use(logger.New())

	// Crie uma instância do serviço de corrida
	corridaService := services.NewCorridaService()

	// Grupo de rotas da API
	api := app.Group("/", logger.New())

	// Rota de Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
	})

	// Configura todas as rotas
	SetupDummyRoutes(api)
	SetupMotoristaRoutes(api)
	SetupCorridaRoutes(api, corridaService)
	NotificacaoCorridaRoutes(api)
}
