package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/gjcms/taxi_service/database"
	"github.com/gjcms/taxi_service/models" // Para cálculos de distância
	"github.com/gjcms/taxi_service/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// LocationUpdateRequest representa o corpo da requisição para atualização de localização
type LocationUpdateRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UpdateDriverLocation lida com a atualização de localização de um motorista
func UpdateDriverLocation(c *fiber.Ctx) error {
	driverIDParam := c.Params("driverID")
	driverID, err := strconv.ParseUint(driverIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid driver ID"})
	}

	var req LocationUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Buscar a corrida em andamento (ou aceita) para este motorista
	var ride models.Ride
	// Procura uma corrida que está aceita ou em progresso e pertence ao motorista
	result := database.DB.Where("driver_id = ? AND (status = ? OR status = ?)", uint(driverID), "accepted", "in_progress").First(&ride)

	if result.Error != nil {
		// Não encontrou corrida ativa, ou houve um erro no DB
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "No active ride for this driver"})
		}
		log.Printf("Error fetching active ride for driver %d: %v", driverID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check active ride"})
	}

	// Lógica para marcar a corrida como "em progresso" quando o motorista se aproxima do passageiro
	if ride.Status == "accepted" {
		distanceToPickup := utils.CalculateDistance(req.Latitude, req.Longitude, ride.OriginLatitude, ride.OriginLongitude)
		if distanceToPickup < 0.1 { // Ex: 100 metros do passageiro
			ride.Status = "in_progress"
			now := time.Now()
			ride.StartedAt = &now // Marca o início real da corrida
			if res := database.DB.Save(&ride); res.Error != nil {
				log.Printf("Error updating ride status to in_progress for ride %d: %v", ride.ID, res.Error)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update ride status"})
			}
			log.Printf("Ride %d status updated to 'in_progress'. Driver %d is at pickup location.", ride.ID, driverID)
		}
	}

	// Lógica para notificação de chegada ao destino (Cenário 1)
	if ride.Status == "in_progress" {
		distanceToDestination := utils.CalculateDistance(req.Latitude, req.Longitude, ride.DestLatitude, ride.DestLongitude)
		if distanceToDestination < 0.05 { // Ex: 50 metros do destino
			// Não mudamos o status para "completed" aqui, apenas notificamos.
			// A finalização da corrida é um evento separado (complete ride).
			log.Printf("Driver %d arrived at destination for ride %d", driverID, ride.ID)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":      "Driver location updated",
				"notification": "Você chegou ao destino", // Sinal para o app cliente exibir
				"ride_status":  ride.Status,              // Para o cliente verificar o status atual
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Location updated, ride status unchanged."})
}

// GetDriverRideHistory retorna o histórico de corridas de um motorista
func GetDriverRideHistory(c *fiber.Ctx) error {
	driverIDParam := c.Params("driverID")
	driverID, err := strconv.ParseUint(driverIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid driver ID"})
	}

	var rides []models.Ride
	// Buscar corridas completas e canceladas para o motorista
	result := database.DB.Where("driver_id = ? AND (status = ? OR status = ?)", uint(driverID), "completed", "cancelled").Order("created_at desc").Find(&rides)
	if result.Error != nil {
		log.Printf("Error fetching ride history for driver %d: %v", driverID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve ride history"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"history": rides})
}
