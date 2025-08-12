package routes

import (
	"github.com/gjcms/taxi_service/controllers"
	"github.com/gofiber/fiber/v2"
)

func DummyRoutes(api fiber.Router) {
	dummy := api.Group("/dummy")

	dummy.Get("/", controllers.ListDummyInfo)
	dummy.Get("/:id", controllers.GetDummyInfo)
	dummy.Post("/", controllers.CreateDummyInfo)
	dummy.Put("/:id", controllers.UpdateDummyInfo)
	dummy.Delete("/:id", controllers.DeleteDummyInfo)
}
