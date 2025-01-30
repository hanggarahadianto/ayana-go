package models

import (
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	ID                     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	WorkerName             string    `json:"worker_name" form:"worker_name"`
	Position               string    `json:"position" form:"position"`
	WeeklyProgressIdWorker uuid.UUID `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"weekly_progress_id"`
	CreatedAt              time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt              time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
