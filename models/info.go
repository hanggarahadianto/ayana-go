package models

import (
	"time"

	"github.com/google/uuid"
)

type Info struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Maps        string    `json:"maps" form:"maps"`
	Start_Price string    `json:"start_price" form:"start_price"`
	Home_ID     string    `gorm:"constraint:OnDelete:CASCADE;"`                                 // Add OnDelete constraint
	NearBy      []NearBy  `gorm:"foreignKey:InfoID;constraint:OnDelete:CASCADE;" json:"nearBy"` // One-to-many relationship with NearBy

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
