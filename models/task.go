package models

import (
	"time"

	"github.com/google/uuid"
)

type Division string

const (
	ITSupport Division = "it_support"
	Marketing Division = "marketing"
	Finance   Division = "finance"
	Legal     Division = "legal"
)

type Task struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string    `json:"name" form:"name"`
	Division Division  `json:"division" form:"division" `
	Note     string    `json:"note" form:"note"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

func (Task) TableName() string {
	return "tasks" // Ensure compatibility with Supabase
}
