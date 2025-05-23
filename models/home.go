package models

import (
	"time"

	"github.com/google/uuid"
)

type Home struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ClusterID  *uuid.UUID `gorm:"type:uuid;index" json:"cluster_id"`
	Type       string     `gorm:"type:varchar(255);not null" json:"type" form:"type"`
	Title      string     `gorm:"type:varchar(255);not null" json:"title" form:"title"`
	Content    string     `gorm:"type:text;not null" json:"content" form:"content"`
	Bathroom   float64    `gorm:"type:bigint;not null"  json:"bathroom" form:"bathroom"`
	Bedroom    float64    `gorm:"type:bigint;not null"  json:"bedroom" form:"bedroom"`
	Square     float64    `gorm:"type:bigint;not null"  json:"square" form:"square"`
	Price      float64    `gorm:"type:bigint;not null" json:"price" form:"price"`
	Quantity   float64    `gorm:"type:bigint;not null"  json:"quantity" form:"quantity"`
	Status     string     `gorm:"type:varchar(255);not null" json:"status"`
	Sequence   int        `gorm:"type:bigint;not null" json:"sequence" form:"sequence"`
	StartPrice float64    `gorm:"type:bigint;not null" json:"start_price" form:"start_price"`
	CreatedAt  time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null" json:"updated_at"`

	Cluster *Cluster `json:"cluster"` // tambah ini

	Images []HomeImage `gorm:"foreignKey:HomeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"images"`
}

type HomeImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	HomeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	ImageURL  string    `gorm:"not null"`
	CreatedAt time.Time
}
