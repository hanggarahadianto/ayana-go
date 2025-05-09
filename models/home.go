package models

import (
	"time"

	"github.com/google/uuid"
)

type Home struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ClusterID  *uuid.UUID `gorm:"type:uuid;index" json:"cluster_id"`
	Title      string     `gorm:"type:varchar(255);not null" json:"title" form:"title"`
	Location   string     `gorm:"type:varchar(255)" json:"location" form:"location"`
	Content    string     `gorm:"type:varchar(255);not null" json:"content" form:"content"`
	Image      string     `gorm:"type:varchar(255);not null" json:"image" form:"image"`
	Bathroom   float64    `gorm:"type:bigint;not null"  json:"bathroom" form:"bathroom"`
	Bedroom    float64    `gorm:"type:bigint;not null"  json:"bedroom" form:"bedroom"`
	Square     float64    `gorm:"type:bigint;not null"  json:"square" form:"square"`
	Price      float64    `gorm:"type:bigint;not null" json:"price" form:"price"`
	Quantity   float64    `gorm:"type:bigint;not null"  json:"quantity" form:"quantity"`
	Status     string     `gorm:"type:varchar(255);not null" json:"status"`
	Sequence   int        `gorm:"type:bigint;not null" json:"sequence" form:"sequence"`
	StartPrice string     `gorm:"type:varchar(255)" json:"start_price" form:"start_price"`
	CreatedAt  time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null" json:"updated_at"`

	NearBies []NearBy `gorm:"foreignKey:HomeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"near_bies"`
}

type NearBy struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `json:"name" form:"name"`
	Distance string    `json:"distance" form:"distance"`
	HomeID   uuid.UUID `gorm:"type:uuid" json:"home_id"`
}

type HomeImage struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	HomeID    string `gorm:"type:uuid;index;not null"`
	ImageURL  string `gorm:"not null"`
	CreatedAt time.Time
}
