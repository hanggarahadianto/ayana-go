package models

import (
	"time"

	"github.com/google/uuid"
)

type WeeklyProgress struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	WeekNumber      string    `json:"week_number" form:"week_number"`
	Percentage      string    `json:"percentage" form:"percentage"`
	Amount_Worker   float64   `json:"amount_worker" form:"amount_worker"`
	Amount_Material float64   `json:"amount_material" form:"amount_material"`
	ProjectID       uuid.UUID `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"project_id"` // Change HomeID type to uuid.UUID
	Project         Project   `gorm:"foreignKey:ProjectID" json:"project"`

	Worker   []Worker   `gorm:"foreignKey:WeekyProgressID" json:"worker"`
	Material []Material `gorm:"foreignKey:MaterialID" json:"material"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
