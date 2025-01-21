package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name   string    `json:"name" form:"name"`
	Email  string    `json:"email" form:"email" gorm:"unique"`
	Phone  string    `json:"phone" form:"phone"`
	HomeID uuid.UUID `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"home_id"` // Change HomeID type to uuid.UUID
	Home   Home      `gorm:"foreignKey:HomeID" json:"home"`                         // Add foreign key relationship with Home

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
