package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectName   string     `gorm:"type:varchar(255);not null" json:"project_name" form:"project_name"`
	ProjectLeader string     `gorm:"type:varchar(255);not null" json:"project_leader" form:"project_leader"`
	Investor      string     `gorm:"type:varchar(255);not null" json:"investor" form:"investor"`
	ProjectTime   string     `gorm:"type:varchar(255);not null" json:"project_time" form:"project_time"`
	TotalCost     int64      `gorm:"type:bigint;not null" json:"total_cost" form:"total_cost"`
	ProjectStart  *time.Time `gorm:"type:timestamp" json:"project_start" form:"project_start"`
	ProjectEnd    *time.Time `gorm:"type:timestamp" json:"project_end" form:"project_end"`
	ProjectStatus *time.Time `gorm:"type:timestamp" json:"project_status" form:"project_status"`
	Note          string     `gorm:"type:varchar(255);not null" json:"note" form:"note"`

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`

	Infos []Info `gorm:"foreignKey:ProjectID" json:"infos"`
}
