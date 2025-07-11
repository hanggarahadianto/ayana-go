// models/presence_rule.go
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TimeTolerant struct {
	Level  int `json:"level" form:"level"`   // Level 1, 2, 3
	Minute int `json:"minute" form:"minute"` // Jumlah menit
}

type PresenceRule struct {
	ID                  uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CompanyID           uuid.UUID     `gorm:"type:uuid;not null" json:"company_id"`
	Day                 string        `gorm:"type:varchar(10);not null" json:"day"`
	IsHoliday           bool          `gorm:"default:false" json:"is_holiday"`
	StartTime           string        `gorm:"type:varchar(8);not null" json:"start_time"` // eg. "07:00"
	EndTime             string        `gorm:"type:varchar(8);not null" json:"end_time"`
	GracePeriodMins     int           `gorm:"default:0" json:"grace_period_mins"`
	ArrivalTolerances   pq.Int32Array `gorm:"type:integer[];default:'{}'" json:"arrival_tolerances"`
	DepartureTolerances pq.Int32Array `gorm:"type:integer[];default:'{}'" json:"departure_tolerances"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
