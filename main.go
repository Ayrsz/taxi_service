package main

import (
    "log"

    "github.com/gjcms/taxi_service/config"
    "github.com/gjcms/taxi_service/database"
    "github.com/gjcms/taxi_service/routes"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors" // <- Adicione esta linha
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    config.LoadConfig()

    database.ConnectDB()

    app := fiber.New()

    // --- Adicione este bloco de código AQUI ---
    app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:5173", // Permite apenas requisições do seu front-end
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))
    // ------------------------------------------

    app.Use(logger.New())

    routes.SetupRoutes(app)

    // Ajuste aqui para usar a porta 3000, conforme seu código
    log.Fatal(app.Listen(":3000")) 
}