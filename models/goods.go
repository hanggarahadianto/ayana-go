package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Decimal float64

// MarshalJSON untuk membulatkan ke 2 desimal
func (d Decimal) MarshalJSON() ([]byte, error) {
	formatted := float64(d)
	// Membulatkan ke 2 desimal
	return json.Marshal(float64(int(formatted*100+0.5)) / 100)
}

type Goods struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GoodsName string    `gorm:"type:varchar(255);not null" json:"good_name" form:"good_name"`
	Status    string    `gorm:"type:varchar(255);not null" json:"status" form:"status"`
	Quantity  Decimal   `gorm:"type:decimal(10,2);not null" json:"quantity" form:"quantity"` // Gunakan tipe Decimal
	CostsDue  float64   `gorm:"type:decimal(10,2)" json:"costs_due" form:"costs_due"`
	Price     float64   `gorm:"type:decimal(10,2)" json:"price" form:"price"`
	Unit      string    `gorm:"type:varchar(255);not null" json:"unit" form:"unit"`
	TotalCost float64   `gorm:"type:decimal(10,2);not null" json:"total_cost" form:"total_cost"`

	CashFlowId uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE;" json:"cash_flow_id"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
