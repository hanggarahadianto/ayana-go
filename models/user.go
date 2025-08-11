package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Username string    `json:"username" form:"username"`
	Password string    `json:"password" form:"password"`
	Role     string    `json:"role" form:"role"`

	Companies []UserCompany `gorm:"foreignKey:UserID" json:"companies"`

	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type RegisterData struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Role            string `json:"role" form:"role"`
}

type LoginData struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
