package middlewares

import (
	"log"
	"strings"

	"github.com/gjcms/taxi_service/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired é um middleware para verificar tokens JWT
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token format is 'Bearer <token>'"})
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Valida o algoritmo de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(config.GlobalConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("JWT Token validation error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse token claims"})
		}

		// Armazena as informações do usuário no contexto do Fiber para handlers posteriores
		c.Locals("user_id", uint(claims["user_id"].(float64))) // JWT claims são float64 para números
		c.Locals("user_email", claims["email"].(string))
		c.Locals("user_role", claims["role"].(string))

		return c.Next() // Continua para o próximo handler ou rota
	}
}
