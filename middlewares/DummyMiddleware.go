package middlewares

import (
	"strings"
	// "taxi-service/models"
    // "taxi-service/services"
	"github.com/gofiber/fiber/v2"
)

func DummyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		c.Locals("dummyToken", token)

		return c.Next()
	}
}
