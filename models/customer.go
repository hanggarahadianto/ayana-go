package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Address       string     `gorm:"type:text;not null" json:"address" form:"address"`
	Phone         string     `gorm:"type:varchar(20);not null" json:"phone"`
	Status        string     `gorm:"type:varchar(100);not null" json:"status" form:"status"`     // contoh: "pending", "deal", etc
	Marketer      string     `gorm:"type:varchar(255);not null" json:"marketer" form:"marketer"` // bisa juga relasi jika marketer entitas sendiri
	Amount        int64      `gorm:"not null" json:"amount"`
	PaymentMethod string     `gorm:"type:varchar(20);not null" json:"payment_method"`
	DateInputed   *time.Time `gorm:"type:timestamp;" json:"date_inputed"` // Tanggal transaksi

	HomeID      *uuid.UUID `gorm:"type:uuid;index" json:"home_id" form:"home_id"`
	ProductUnit string     `gorm:"type:varchar(20);not null" json:"product_unit" form:"product_unit"`
	BankName    string     `gorm:"type:varchar(20);" json:"bank_name" form:"bank_name"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Home *Home `json:"home"` // Relasi ke model Home
}
