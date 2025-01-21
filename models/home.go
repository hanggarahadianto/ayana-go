package models

import (
	"time"

	"github.com/google/uuid"
)

type Home struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title    string    `json:"title" form:"title"`
	Content  string    `json:"content" form:"content"`
	Image    string    `json:"image" form:"image"`
	Address  string    `json:"address" form:"address"`
	Bathroom string    `json:"bathroom" form:"bathroom"`
	Bedroom  string    `json:"bedroom" form:"bedroom"`
	Square   string    `json:"square" form:"square"`
	Price    float64   `json:"price" form:"price"`
	Quantity int       `json:"quantity" form:"quantity"`
	Status   string    `json:"status"`
	// Remove the direct Info relationship to avoid circular references
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
