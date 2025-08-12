package models

import "time"

// User representa um usuário no sistema (motorista ou passageiro)
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // Senha hashed, não serializada para JSON
	Role      string    `json:"role"`              // ex: "driver", "passenger"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Outros campos como Nome, Telefone, Avatar, etc.
}