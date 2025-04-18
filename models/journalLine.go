package models

import (
	"time"

	"github.com/google/uuid"
)

type JournalLine struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	JournalID   uuid.UUID `gorm:"type:uuid;not null;onDelete:CASCADE" json:"journal_id"` // Menambahkan onDelete: "CASCADE"
	AccountID   uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	Debit       int64     `json:"debit"`
	Credit      int64     `json:"credit"`
	Description string    `gorm:"type:varchar(255)" json:"description"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Account Account `gorm:"foreignKey:AccountID" json:"account"`
}
