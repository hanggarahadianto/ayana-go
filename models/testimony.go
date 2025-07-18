package models

import (
	"time"

	"github.com/google/uuid"
)

type Testimony struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CustomerID *uuid.UUID `gorm:"type:uuid;index" json:"customer_id" form:"customer_id"` // ganti dari HomeID
	CompanyID  uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	Rating     int        `gorm:"type:int;check:rating >= 1 AND rating <= 5" json:"rating"` // batas rating 1-5
	Note       string     `gorm:"type:text" json:"note"`                                    // komentar pengguna

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer"` // optional: eager load customer data
}
