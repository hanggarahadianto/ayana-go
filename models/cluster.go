package models

import (
	"time"

	"github.com/google/uuid"
)

type Cluster struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Location   string    `gorm:"type:varchar(255)" json:"location" form:"location"`
	Image      string    `gorm:"type:varchar(255);not null" json:"image" form:"image"`
	Address    string    `gorm:"type:varchar(255);not null" json:"address" form:"address"`
	Bathroom   float64   `gorm:"type:bigint;not null"  json:"bathroom" form:"bathroom"`
	Bedroom    float64   `gorm:"type:bigint;not null"  json:"bedroom" form:"bedroom"`
	Square     float64   `gorm:"type:bigint;not null"  json:"square" form:"square"`
	Price      float64   `gorm:"type:bigint;not null" json:"price" form:"price"`
	Quantity   float64   `gorm:"type:bigint;not null"  json:"quantity" form:"quantity"`
	Status     string    `gorm:"type:varchar(255);not null" json:"status"`
	Sequence   int       `gorm:"type:bigint;not null" json:"sequence" form:"sequence"`
	Maps       string    `gorm:"type:varchar(255)" json:"maps" form:"maps"`
	StartPrice string    `gorm:"type:varchar(255)" json:"start_price" form:"start_price"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null" json:"updated_at"`
}
