package models

import (
	"time"

	"github.com/google/uuid"
)

type WeeklyProgress struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	WeekNumber     string    `json:"week_number" form:"week_number"`
	Percentage     int64     `gorm:"type:bigint;not null" json:"percentage" form:"percentage"`
	AmountWorker   int64     `gorm:"type:bigint;not null" json:"amount_worker" form:"amount_worker"`
	AmountMaterial int64     `gorm:"type:bigint;not null" json:"amount_material" form:"amount_material"`

	ProjectID uuid.UUID `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"project_id"`

	Material []Material `gorm:"foreignKey:WeeklyProgressIdMaterial" json:"material"`
	Worker   []Worker   `gorm:"foreignKey:WeeklyProgressIdWorker" json:"worker"`
	Note     string     `gorm:"type:varchar(255);not null" json:"note" form:"note"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
