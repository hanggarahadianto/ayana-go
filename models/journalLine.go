package models

import (
	"time"

	"github.com/google/uuid"
)

type JournalLine struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	JournalID   uuid.UUID `gorm:"type:uuid;not null" json:"journal_id"`          // Foreign key untuk JournalEntry
	AccountCode string    `gorm:"type:varchar(50);not null" json:"account_code"` // Kode akun
	Debit       int64     `gorm:"type:bigint;default:0" json:"debit"`            // Jumlah debit
	Credit      int64     `gorm:"type:bigint;default:0" json:"credit"`           // Jumlah kredit
	Description string    `gorm:"type:varchar(255)" json:"description"`          // Deskripsi baris jurnal
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
