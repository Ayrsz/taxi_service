package routes

import (
	"taxi_service/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupDummyRoutes(api fiber.Router) {
	// Rotas de dummy users
	dummyUsers := api.Group("/dummy-users")
	dummyUsers.Get("/", controllers.ListDummyInfo)
	dummyUsers.Get("/:id", controllers.GetDummyInfo)
	dummyUsers.Post("/", controllers.CreateDummyInfo)
	dummyUsers.Put("/:id", controllers.UpdateDummyInfo)
	dummyUsers.Delete("/:id", controllers.DeleteDummyInfo)
}
