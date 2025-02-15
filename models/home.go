package models

import (
	"time"

	"github.com/google/uuid"
)

type Home struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`

	Title     string    `gorm:"type:varchar(255);not null" json:"title" form:"title"`
	Location  string    `gorm:"type:varchar(255);not null" json:"location" form:"location"`
	Content   string    `gorm:"type:varchar(255);not null" json:"content" form:"content"`
	Image     string    `gorm:"type:varchar(255);not null" json:"image" form:"image"`
	Address   string    `gorm:"type:varchar(255);not null" json:"address" form:"address"`
	Bathroom  string    `gorm:"type:varchar(255);not null" json:"bathroom" form:"bathroom"`
	Bedroom   string    `gorm:"type:varchar(255);not null" json:"bedroom" form:"bedroom"`
	Square    string    `gorm:"type:varchar(255);not null" json:"square" form:"square"`
	Price     float64   `gorm:"type:bigint;not null" json:"price" form:"price"`
	Quantity  int       `gorm:"type:bigint;not null" json:"quantity" form:"quantity"`
	Status    string    `gorm:"type:varchar(255);not null" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
