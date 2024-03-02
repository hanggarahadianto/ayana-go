package models

import (
	"time"

	"github.com/google/uuid"
)

type Info struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Office_Phone string    `json:"officePhone" form:"officePhone"`
	Email        string    `json:"officeEmail" form:"officeEmail"`
	Address      string    `json:"address" form:"address"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
