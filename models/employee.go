package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID                   uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name                 string     `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Address              string     `gorm:"type:text;not null" json:"address" form:"address"`
	Phone                string     `gorm:"type:varchar(20);not null" json:"phone"`
	DateBirth            *time.Time `gorm:"type:timestamp" json:"date_birth" form:"date_birth"`
	MaritalStatus        string     `gorm:"type:varchar(50);not null" json:"marital_status" form:"marital_status"`  // contoh: "single", "married", etc
	EmployeeEducation    string     `gorm:"type:text;not null" json:"employee_education" form:"employee_education"` // pendidikan terakhir
	Department           string     `gorm:"type:text;not null" json:"department" form:"address"`
	Gender               string     `gorm:"type:text;not null" json:"gender" form:"gender"`
	Religion             string     `gorm:"type:text;not null" json:"religion" form:"religion"` // contoh: "Islam", "Kristen", etc
	Position             string     `gorm:"type:text;not null" json:"position" form:"position"`
	EmployeeStatus       string     `gorm:"type:varchar(100);not null" json:"employee_status" form:"employee_status"`               // contoh: "active", "inactive", etc
	EmployeeContractType string     `gorm:"type:varchar(100);not null" json:"employee_contract_type" form:"employee_contract_type"` // contoh: "permanent", "contract", etc
	CompanyID            uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
