package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Code        int16     `gorm:"type:varchar(20);unique;not null" json:"code"` // ex: 101, 201
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`       // ex: Kas, Gaji
	Type        string    `gorm:"type:varchar(50);not null" json:"type"`        // Asset, Expense, Liability, Revenue, Equity
	Category    string    `gorm:"type:varchar(100)" json:"category"`            // Gaji, Operasional, Piutang, dll
	Description string    `gorm:"type:varchar(255)" json:"description"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"` // Foreign key

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
