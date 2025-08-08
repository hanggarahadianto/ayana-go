package models

import (
	"time"

	"github.com/google/uuid"
)

type UserCompany struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`

	// Relasi opsional jika mau preload data user atau company
	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Company Company `gorm:"foreignKey:CompanyID" json:"company"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
