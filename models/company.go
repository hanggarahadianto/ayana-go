package models

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title       string    `json:"title" form:"title"`
	CompanyCode string    `gorm:"unique;not null" json:"company_code" form:"company_code"`
	Color       string    `json:"color" form:"color"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
