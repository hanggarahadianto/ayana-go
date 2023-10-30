package models

import (
	"time"

	"github.com/google/uuid"
)

type Home struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title    string    `json:"title" form:"title"`
	Content  string    `json:"content" form:"content"`
	Image    string    `json:"image" form:"image"`
	Address  string    `json:"address" form:"address"`
	Bathroom string    `json:"bathroom" form:"bathroom"`
	Bedroom  string    `json:"bedroom" form:"bedroom"`
	Square   string    `json:"square" form:"square"`

	Reservation []Reservation `gorm:"Foreignkey:Home_ID;association_foreignkey:ID;" json:"reservation"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type UpdateHome struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	// Image    string `json:"image" form:"image"`
	Address  string `json:"address" form:"address"`
	Bathroom string `json:"bathroom" form:"bathroom"`
	Bedroom  string `json:"bedroom" form:"bedroom"`
	Square   string `json:"square" form:"square"`
}
