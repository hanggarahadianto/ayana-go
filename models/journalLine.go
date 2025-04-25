package models

import (
	"time"

	"github.com/google/uuid"
)

type JournalLine struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	JournalID   uuid.UUID `gorm:"type:uuid;not null;onDelete:CASCADE" json:"journal_id"` // Menambahkan onDelete: "CASCADE"
	AccountID   uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	AccountName string    `gorm:"type:varchar(100)" json:"account_name"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"` // Tambahkan kolom company_id

	Debit             int64  `json:"debit"`
	CreditAccountName string `gorm:"type:varchar(100)" json:"credit_account_name"` // Nama akun kredit

	Credit int64 `json:"credit"`

	Description     string          `gorm:"type:varchar(255)" json:"description"`
	TransactionType TransactionType `gorm:"type:varchar(50)" json:"transaction_type"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Account Account `gorm:"foreignKey:AccountID" json:"account"`
	// Company Company      `gorm:"foreignkey:CompanyID" json:"company"`
	Journal JournalEntry `gorm:"foreignKey:JournalID;references:ID" json:"journal"`
}
