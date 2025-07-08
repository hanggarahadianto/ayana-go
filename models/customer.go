package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name    string    `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Address string    `gorm:"type:text;not null" json:"address" form:"address"`
	Phone   string    `gorm:"type:varchar(20);not null" json:"phone"`

	Status       string    `gorm:"type:varchar(100);not null" json:"status" form:"status"`   // contoh: "pending", "deal", etc
	MarketerID   uuid.UUID `gorm:"type:uuid;not null" json:"marketer_id" form:"marketer_id"` // FK ke Employee
	Marketer     *Employee `gorm:"foreignKey:MarketerID" json:"marketer"`                    // ⬅️ ubah json tag agar tidak konflik
	MarketerName string    `gorm:"type:varchar(100);not null" json:"marketer_name" form:"marketer_name"`

	Amount        int64      `gorm:"not null" json:"amount"`
	PaymentMethod string     `gorm:"type:varchar(20);not null" json:"payment_method"`
	DateInputed   *time.Time `gorm:"type:timestamp;" json:"date_inputed"` // Tanggal transaksi

	HomeID      *uuid.UUID `gorm:"type:uuid;index" json:"home_id" form:"home_id"`
	ProductUnit string     `gorm:"type:varchar(20);not null" json:"product_unit" form:"product_unit"`
	BankName    string     `gorm:"type:varchar(20);" json:"bank_name" form:"bank_name"`

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Home *Home `json:"home"` // Relasi ke model Home

}
