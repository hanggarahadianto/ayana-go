package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Address   string     `gorm:"type:text;not null" json:"address" form:"address"`
	Phone     string     `gorm:"type:varchar(20);not null" json:"phone"`
	Status    string     `gorm:"type:varchar(100);not null" json:"status" form:"status"`     // contoh: "pending", "deal", etc
	Marketer  string     `gorm:"type:varchar(255);not null" json:"marketer" form:"marketer"` // bisa juga relasi jika marketer entitas sendiri
	HomeID    *uuid.UUID `gorm:"type:uuid;index" json:"home_id" form:"home_id"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`

	Home *Home `json:"home"` // Relasi ke model Home
}
