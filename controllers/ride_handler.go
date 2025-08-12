package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/gjcms/taxi_service/utils"
	"github.com/gjcms/taxi_service/database"
	"github.com/gjcms/taxi_service/models" // Para cálculos de ETA
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RequestRidePayload representa o corpo da requisição para solicitar uma corrida (simplificado)
type RequestRidePayload struct {
	PassengerID uint    `json:"passenger_id"`
	OriginLat   float64 `json:"origin_latitude"`
	OriginLon   float64 `json:"origin_longitude"`
	DestLat     float64 `json:"dest_latitude"`
	DestLon     float64 `json:"dest_longitude"`
}

// RequestRide lida com a solicitação de uma nova corrida por um passageiro
func RequestRide(c *fiber.Ctx) error {
	var req RequestRidePayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Calcular distância estimada e valor (simplificado)
	estimatedDistance := utils.CalculateDistance(req.OriginLat, req.OriginLon, req.DestLat, req.DestLon)
	estimatedValue := estimatedDistance * 2.5 // Exemplo: R$2.50 por km

	ride := models.Ride{
		PassengerID:         req.PassengerID,
		OriginLatitude:      req.OriginLat,
		OriginLongitude:     req.OriginLon,
		DestLatitude:        req.DestLat,
		DestLongitude:       req.DestLon,
		Status:              "pending", // Corrida pendente esperando um motorista
		EstimatedDistanceKM: estimatedDistance,
		EstimatedValue:      estimatedValue,
		CreatedAt:           time.Now(),
	}

	if result := database.DB.Create(&ride); result.Error != nil {
		log.Printf("Error creating ride: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not request ride"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Ride requested successfully", "ride": ride})
}

// AcceptRide lida com a aceitação de uma corrida por um motorista (Cenário 5)
func AcceptRide(c *fiber.Ctx) error {
	rideIDParam := c.Params("id")
	rideID, err := strconv.ParseUint(rideIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ride ID"})
	}
	var ride models.Ride
	// Verifique se a corrida existe e está pendente
	if result := database.DB.Where("status = ?", "pending").First(&ride, uint(rideID)); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ride not found or not pending"})
		}
		log.Printf("Error fetching ride %d: %v", rideID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve ride"})
	}
	driverID := uint(1) // Assumindo que o motorista João (ID 1) está logado para os testes

	ride.DriverID = driverID
	ride.Status = "accepted"
	now := time.Now()
	ride.AcceptedAt = &now

	// Calcular ETA até o local de embarque (seria a localização ATUAL do motorista até o ponto de embarque)
	
	ride.ETA = utils.CalculateETA(ride.OriginLatitude, ride.OriginLongitude, 0.5) // Ex: 0.5 km/min de velocidade média

	if result := database.DB.Save(&ride); result.Error != nil {
		log.Printf("Error accepting ride %d: %v", rideID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not accept ride"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Ride accepted successfully",
		"ride":         ride,
		"notification": "Tempo estimado de chegada: " + ride.ETA, // Sinal para o app cliente exibir
	})
}

// CompleteRideRequest representa o corpo da requisição para finalizar uma corrida
type CompleteRideRequest struct {
	ActualDistanceKM float64 `json:"actual_distance_km"`
	ActualValue      float64 `json:"actual_value"`
}

// CompleteRide lida com a finalização de uma corrida (Cenário 3)
func CompleteRide(c *fiber.Ctx) error {
	rideIDParam := c.Params("id")
	rideID, err := strconv.ParseUint(rideIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ride ID"})
	}

	var req CompleteRideRequest 
	if err := c.BodyParser(&req); err != nil {
		// Se o parse do corpo falhar, significa que os valores reais não foram enviados
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body for completion"})
	}

	var ride models.Ride
	// Verifique se a corrida existe e está em progresso
	if result := database.DB.Where("status = ?", "in_progress").First(&ride, uint(rideID)); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ride not found or not in progress"})
		}
		log.Printf("Error fetching ride %d: %v", rideID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve ride"})
	}
	ride.Status = "completed"
	now := time.Now()
	ride.CompletedAt = &now
	// Usar os valores do corpo da requisição
	ride.ActualDistanceKM = req.ActualDistanceKM
	ride.ActualValue = req.ActualValue

	if result := database.DB.Save(&ride); result.Error != nil {
		log.Printf("Error completing ride %d: %v", rideID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not complete ride"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Ride completed successfully", "ride": ride})
}

// CancelRideRequest representa o corpo da requisição para cancelar corrida (se houver motivo)
type CancelRideRequest struct {
	Reason string `json:"reason"`
}

// CancelRide lida com o cancelamento de uma corrida (Cenário 2)
func CancelRide(c *fiber.Ctx) error {
	rideIDParam := c.Params("id")
	rideID, err := strconv.ParseUint(rideIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ride ID"})
	}

	var req CancelRideRequest
	// Faz o parse, 
	_ = c.BodyParser(&req)

	var ride models.Ride
	// Verifique se a corrida existe e está em um status que pode ser cancelado
	if result := database.DB.Where("status IN ?", []string{"pending", "accepted", "in_progress"}).First(&ride, uint(rideID)); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ride not found or cannot be cancelled in current status"})
		}
		log.Printf("Error fetching ride %d for cancellation: %v", rideID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve ride"})
	}

	// Lógica para exibir a notificação de confirmação antes do cancelamento real
	// envia-se uma primeira requisição para "cancelar" e espera esta resposta de confirmação.
	// Se se confirmar,é enviado outra requisição 
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Confirmation required for cancellation",
		"notification": "Tem certeza que deseja cancelar a corrida? Cancelamentos frequentes podem impactar sua avaliação.",
		"options":      []string{"Sim, quero cancelar", "Não, continuar com a corrida"},
		"ride_id":      ride.ID, // Retorna o ID da corrida para usar na confirmação
	})

	// Se você quiser implementar o cancelamento final aqui, seria algo assim:
	/*
		ride.Status = "cancelled"
		now := time.Now()
		ride.CancelledAt = &now
		ride.CancellationReason = req.Reason // Usa o motivo da requisição

		if result := database.DB.Save(&ride); result.Error != nil {
			log.Printf("Error cancelling ride %d: %v", rideID, result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not cancel ride"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Ride cancelled successfully", "ride": ride})
	*/
}