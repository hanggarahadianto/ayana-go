package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectStatus string

const (
	StatusReady      ProjectStatus = "ready"
	StatusOnProgress ProjectStatus = "on_progress"
	StatusDone       ProjectStatus = "done"
)

type Project struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectName   string     `gorm:"type:varchar(255);not null" json:"project_name" form:"project_name"`
	Location      string     `gorm:"type:varchar(255);not null" json:"location" form:"location"`
	Type          string     `gorm:"type:varchar(255);not null" json:"type" form:"type"`
	Unit          string     `gorm:"type:varchar(255);not null" json:"unit" form:"unit"`
	ProjectLeader string     `gorm:"type:varchar(255);not null" json:"project_leader" form:"project_leader"`
	Investor      string     `gorm:"type:varchar(255);not null" json:"investor" form:"investor"`
	ProjectTime   string     `gorm:"type:varchar(255);not null" json:"project_time" form:"project_time"`
	TotalCost     int64      `gorm:"type:bigint;not null" json:"total_cost" form:"total_cost"`
	ProjectStart  *time.Time `gorm:"type:timestamp" json:"project_start" form:"project_start"`
	ProjectEnd    *time.Time `gorm:"type:timestamp" json:"project_end" form:"project_end"`

	ProjectFinished *time.Time `gorm:"type:timestamp" json:"project_finished" form:"project_finished"`

	ProjectStatus ProjectStatus `gorm:"type:varchar(20);not null;check:project_status IN ('ready','on_progress','done')" json:"project_status" form:"project_status"`
	Note          string        `gorm:"type:varchar(255);not null" json:"note" form:"note"`

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()" json:"updated_at"`
}

type ProjectWithStatus struct {
	Project
	StatusText   string `json:"status_text"`
	SisaWaktu    string `json:"sisa_waktu,omitempty"`
	Color        string `json:"color"`
	IsOnTime     *bool  `json:"is_on_time,omitempty"`
	DelayDays    int    `json:"delay_days,omitempty"`
	FinishStatus string `json:"finish_status,omitempty"` // ⬅️ baru
}
