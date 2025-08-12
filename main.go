package main

import (
    "github.com/gofiber/fiber/v2"
    "taxi-service/routes"
    "taxi-service/services"
)

func main() {
	// Carrega as corridas do JSON
	services.CarregarCorridasDoArquivo()

	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":3000")
}
