package controllers

import (
	"time"

	"github.com/gjcms/taxi_service/config"
	"github.com/gjcms/taxi_service/database"
	"github.com/gjcms/taxi_service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // Instale com "go get github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"   // Instale com "go get golang.org/x/crypto/bcrypt"
)

// RegisterUserRequest representa o corpo da requisição para registro de usuário
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=driver passenger"`
}

// LoginRequest representa o corpo da requisição para login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterUser lida com a criação de um novo usuário
func RegisterUser(c *fiber.Ctx) error {
	var req RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	// Retorna uma resposta sem a senha
	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully", "user": user})
}

// Login lida com a autenticação do usuário
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compara a senha fornecida com a senha hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Cria o token JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expira em 24 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GlobalConfig.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful", "token": t})
}
