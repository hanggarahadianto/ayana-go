package models

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	MaterialName string    `gorm:"type:varchar(255);not null" json:"material_name" form:"material_name"`
	Quantity     float64   `gorm:"type:varchar(255);not null" json:"quantity" form:"quantity"`
	TotalCost    float64   `gorm:"type:varchar(255);not null" json:"total_cost" form:"total_cost"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
