package models

import "github.com/google/uuid"

type Worker struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	WorkerName string    `json:"project_name" form:"project_name"`
	Position   uuid.UUID `gorm:"type:position" json:"position"`
}

func (Worker) TableName() string {
	return "worker" // Ensure compatibility with Supabase
}
