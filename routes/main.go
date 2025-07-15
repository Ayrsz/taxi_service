package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes configura todas as rotas da aplicação.
func SetupRoutes(app *fiber.App) {
    // Cria um grupo de rotas na raiz '/' e adiciona um logger
    // para vermos as requisições no console.
    api := app.Group("/", logger.New())

    // Rota de "health check" para saber se a API está no ar.
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
    })

    // Chama a função que configura as rotas de corrida.
    SetupCorridaRoutes(api)
}
