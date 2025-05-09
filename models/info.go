package models

import (
	"time"

	"github.com/google/uuid"
)

type Info struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`

	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	// Relasi ke Project
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
