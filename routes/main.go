package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes configura todas as rotas da aplicação.
func SetupRoutes(app *fiber.App) {
	api := app.Group("/", logger.New())

    // Rota de "health check" para saber se a API está no ar.
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
    })

	SetupDummyRoutes(api)
}
