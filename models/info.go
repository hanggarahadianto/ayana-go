package models

import (
	"time"

	"github.com/google/uuid"
)

type Info struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Maps       string    `json:"maps"`
	StartPrice string    `json:"start_price"`
	HomeID     uuid.UUID `gorm:"type:uuid" json:"home_id"`
	Home       Home      `gorm:"foreignKey:HomeID" json:"-"` // Add this line
	NearBy     []NearBy  `gorm:"foreignKey:InfoID" json:"near_by"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}
