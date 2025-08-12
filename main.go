package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // 1. Importe o middleware de CORS
	"taxi-service/routes"
	"taxi-service/services"
)

func main() {
	// Carrega as corridas do JSON
	services.CarregarCorridasDoArquivo()

	app := fiber.New()

	// 2. Adicione o middleware de CORS aqui
	// A configuração foi tornada mais permissiva para depuração,
	// permitindo qualquer origem e métodos POST.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Permite qualquer origem (bom para depuração)
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.SetupRoutes(app)

	app.Listen(":3000")
}