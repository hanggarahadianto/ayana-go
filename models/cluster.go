package models

import (
	"time"

	"github.com/google/uuid"
)

type Cluster struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Location  string    `gorm:"type:varchar(255)" json:"location" form:"location"`
	Square    float64   `gorm:"type:bigint;not null"  json:"square" form:"square"`
	Price     float64   `gorm:"type:bigint;not null" json:"price" form:"price"`
	Quantity  float64   `gorm:"type:bigint;not null"  json:"quantity" form:"quantity"`
	Status    string    `gorm:"type:varchar(255);not null" json:"status"`
	Sequence  int       `gorm:"type:bigint;not null" json:"sequence" form:"sequence"`
	Maps      string    `gorm:"type:varchar(255)" json:"maps" form:"maps"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Homes []Home `gorm:"foreignKey:ClusterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"homes"`
}
