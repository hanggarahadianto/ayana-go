package models

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID                       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	MaterialName             string    `gorm:"type:varchar(255);not null" json:"material_name" form:"material_name"`
	Quantity                 int       `gorm:"not null" json:"quantity" form:"quantity"`
	TotalCost                float64   `gorm:"type:decimal(10,2);not null" json:"total_cost" form:"total_cost"`
	WeeklyProgressIdMaterial uuid.UUID `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"weekly_progress_id"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
