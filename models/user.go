package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username string    `json:"username" form:"username"`
	Password string    `json:"password" form:"password"`
	Role     string    `json:"role" form:"role"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type RegisterData struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Role            string `json:"role" form:"role"`
}

type LoginData struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
