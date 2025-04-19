package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string
type TransactionType string

const (
	StatusDraft     Status = "draft"
	StatusApproved  Status = "approved"
	StatusPaid      Status = "paid"
	StatusUnpaid    Status = "unpaid"
	StatusCancelled Status = "cancelled"
)
const (
	PayinType  TransactionType = "payin"
	PayoutType TransactionType = "payout"
)

type JournalEntry struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Invoice               string    `gorm:"type:varchar(100)" json:"invoice"`
	Description           string    `gorm:"type:text" json:"description"`
	TransactionCategoryID uuid.UUID `gorm:"type:uuid" json:"transaction_category_id"`
	Amount                int64     `gorm:"not null" json:"amount"`
	Partner               string    `gorm:"type:text;not null" json:"partner"`

	TransactionType TransactionType `gorm:"type:varchar(50)" json:"transaction_type"`
	Status          Status          // e.g., "draft", "posted"
	CompanyID       uuid.UUID       `gorm:"type:uuid;not null" json:"company_id"`

	DateInputed *time.Time `gorm:"type:timestamp;not null" json:"date_inputed"` // Tanggal transaksi
	DueDate     *time.Time `gorm:"type:timestamp" json:"due_date,omitempty"`    // nullable, tergantung jenis transaksi
	IsRepaid    bool       `json:"is_repaid"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Lines []JournalLine `gorm:"foreignKey:JournalID" json:"lines"`

	TransactionCategory TransactionCategory `gorm:"foreignKey:TransactionCategoryID" json:"transaction_category"` // optional relasi
}
