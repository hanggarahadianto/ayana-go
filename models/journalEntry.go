package models

import (
	"time"

	"github.com/google/uuid"
)

type JournalEntry struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Date        *time.Time `gorm:"type:timestamp;not null" json:"date"` // Tanggal transaksi
	Description string     `gorm:"type:varchar(255);not null" json:"description"`
	Invoice     string     `gorm:"type:varchar(100);unique" json:"invoice"`
	Category    string     `gorm:"type:varchar(100)" json:"category"`              // misalnya: payroll, equipment, dsb
	Status      string     `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, posted, canceled

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Lines []JournalLine `gorm:"foreignKey:JournalID" json:"lines"`
}
