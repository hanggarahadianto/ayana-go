package models

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MaterialName string    `gorm:"type:varchar(255);not null" json:"material_name"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	Unit         string    `gorm:"type:varchar(255);not null" json:"unit"`
	Price        int64     `gorm:"type:bigint;not null" json:"price"`
	TotalCost    float64   `gorm:"type:decimal(10,2);not null" json:"total_cost"`

	// Explicit foreign key reference
	WeeklyProgressIdMaterial uuid.UUID      `gorm:"type:uuid;not null" json:"weekly_progress_id"`
	WeeklyProgress           WeeklyProgress `gorm:"foreignKey:WeeklyProgressIdMaterial;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
