package main

import (
    "github.com/gofiber/fiber/v2"
    "taxi_service/routes" // nome do módulo + pasta
    "taxi_service/services"
)

func main() {
	// Carrega as corridas do JSON
	services.CarregarCorridasDoArquivo()

	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":8080")
}
