package models

import (
	"time"

	"github.com/google/uuid"
)

// TransactionCategory adalah model untuk kategori transaksi
type TransactionCategory struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`  // Nama kategori transaksi, misal: "Investor", "Kreditur"
	Status           string    `gorm:"type:varchar(10);not null" json:"status"` // Nilai: "paid" atau "unpaid"
	TransactionLabel string    `gorm:"type:varchar(50);not null" json:"transaction_label"`

	DebitAccountID    uuid.UUID `gorm:"type:uuid;not null" json:"debit_account_id"`        // Foreign key untuk akun debit
	DebitAccountType  string    `gorm:"type:varchar(100)" json:"debit_account_type"`       // Nama akun debit
	CreditAccountID   uuid.UUID `gorm:"type:uuid;not null" json:"credit_account_id"`       // Foreign key untuk akun kredit
	CreditAccountType string    `gorm:"type:varchar(100)" json:"credit_account_type"`      // Nama akun kredit
	DebitCategory     string    `gorm:"type:varchar(100)" json:"debit_category"`           // Kategori umum untuk transaksi (misal: "Pembayaran", "Penerimaan")
	CreditCategory    string    `gorm:"type:varchar(100)" json:"credit_category"`          // Kategori umum untuk transaksi (misal: "Pembayaran", "Penerimaan")
	TransactionType   string    `gorm:"type:varchar(50);not null" json:"transaction_type"` // payout / payin
	Description       string    `gorm:"type:varchar(255)" json:"description"`              // Deskripsi transaksi
	CompanyID         uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`              // Foreign key untuk perusahaan
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	DebitAccount  Account `gorm:"foreignKey:DebitAccountID" json:"debit_account"`   // Relasi dengan Account untuk akun debit
	CreditAccount Account `gorm:"foreignKey:CreditAccountID" json:"credit_account"` // Relasi dengan Account untuk akun kredit
}
