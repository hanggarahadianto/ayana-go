package models

import (
	"time"

	"github.com/google/uuid"
)

type Goods struct {
	ID                 uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GoodsName          string     `gorm:"type:varchar(255);not null" json:"material_name" form:"material_name"`
	Status             string     `gorm:"type:varchar(255);not null" json:"status" form:"status"`
	Quantity           int        `gorm:"not null" json:"quantity" form:"quantity"`
	TotalCost          float64    `gorm:"type:decimal(10,2);not null" json:"total_cost" form:"total_cost"`
	GoodPurchaseDate   *time.Time `gorm:"type:timestamp" json:"good_purcase_date" form:"good_purchase_date"`
	GoodSettlementDate *time.Time `gorm:"type:timestamp" json:"good_settlement_date" form:"good_settlement_date"`
	CashFlowId         uuid.UUID  `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"cash_flow_id"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
