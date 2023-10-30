package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name    string    `json:"name" form:"name"`
	Phone   string    `json:"phone" form:"phone"`
	Date    string    `json:"date" form:"date"`
	Time    string    `json:"time" form:"time"`
	Home_ID string    `gorm:"column:home_id"  json:"home_id"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
