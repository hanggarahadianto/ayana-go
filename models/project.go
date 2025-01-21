package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectName   string     `gorm:"type:varchar(255);not null" json:"project_name" form:"project_name"`
	ProjectLeader string     `gorm:"type:varchar(255);not null" json:"project_leader" form:"project_leader"`
	ProjectTime   string     `gorm:"type:varchar(255);not null" json:"project_time" form:"project_time"`
	TotalCost     float64    `gorm:"type:decimal(10,2);not null" json:"total_cost" form:"total_cost"`
	ProjectStart  *time.Time `gorm:"type:timestamp" json:"project_start" form:"project_start"`
	ProjectEnd    *time.Time `gorm:"type:timestamp" json:"project_end" form:"project_end"`
	Note          string     `gorm:"type:varchar(255);not null" json:"note" form:"note"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`
}
