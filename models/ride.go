package models

import "time"

// Ride representa uma corrida de táxi
type Ride struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	DriverID           uint      `json:"driver_id"` // ID do motorista associado (Chave estrangeira para User)
	PassengerID        uint      `json:"passenger_id"` // ID do passageiro (Chave estrangeira para User)
	OriginLatitude     float64   `json:"origin_latitude"`
	OriginLongitude    float64   `json:"origin_longitude"`
	DestLatitude       float64   `json:"dest_latitude"`
	DestLongitude      float64   `json:"dest_longitude"`
	Status             string    `json:"status"` // ex: "pending", "accepted", "in_progress", "completed", "cancelled"
	EstimatedDistanceKM float64   `json:"estimated_distance_km"`
	EstimatedValue     float64   `json:"estimated_value"`
	ActualDistanceKM   float64   `json:"actual_distance_km"`
	ActualValue        float64   `json:"actual_value"`
	AcceptedAt         *time.Time `json:"accepted_at"` // Ponteiro para permitir valor nulo no DB
	StartedAt          *time.Time `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancellationReason string    `json:"cancellation_reason"` // Motivo do cancelamento
	ETA                string    `json:"eta"`                 // Estimated Time of Arrival to pickup location (para o motorista)
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}