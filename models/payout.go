package models

import (
	"time"

	"github.com/google/uuid"
)

type Payout struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Invoice     string     `json:"invoice" form:"invoice"`
	Nominal     int64      `gorm:"type:bigint;not null" json:"nominal" form:"nominal"`
	DateInputed *time.Time `gorm:"type:timestamp" json:"date_inputed" form:"date_inputed"`
	Note        string     `gorm:"type:varchar(255);not null" json:"note" form:"note"`

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"` // Foreign key ke Company

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
