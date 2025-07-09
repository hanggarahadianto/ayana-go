// models/presence.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Presence struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EmployeeID uuid.UUID `gorm:"type:uuid;not null" json:"employee_id"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID"` // untuk preload

	ScanDate time.Time `json:"scan_date"`
	ScanTime string    `json:"scan_time"`
	RawDate  string    `json:"raw_date"`

	CreatedAt time.Time
	UpdatedAt time.Time

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"` // Foreign key

}
