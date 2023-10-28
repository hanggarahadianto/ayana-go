package models

import (
	"github.com/google/uuid"
)

type Image struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title string    `json:"title" form:"title"`
	Image string    `json:"image" form:"image"`
}
