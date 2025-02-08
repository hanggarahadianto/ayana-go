package models

import (
	"time"

	"github.com/google/uuid"
)

type CashFlow struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	WeekNumber  string    `json:"week_number" form:"week_number"`
	CashIn      int64     `gorm:"type:bigint;not null" json:"cash_in" form:"cash_in"`
	CashOut     int64     `gorm:"type:bigint;not null" json:"cash_out" form:"cash_out"`
	Outstanding int64     `gorm:"type:bigint;not null" json:"outstanding" form:"outstanding"`

	ProjectID uuid.UUID `gorm:"type:uuid" json:"project_id"`

	Good []Goods `gorm:"foreignKey:CashFlowId;constraint:OnDelete:CASCADE;" json:"good"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
