package models

import "github.com/google/uuid"

type NearBy struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name     string    `json:"name" form:"name"`
	Distance string    `json:"distance" form:"distance"`
	InfoID   uuid.UUID `gorm:"type:uuid;constraint:OnDelete:CASCADE;" json:"info_id"`
}
