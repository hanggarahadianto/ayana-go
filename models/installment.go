package models

import (
	"time"

	"github.com/google/uuid"
)

type Installment struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	JournalID uuid.UUID    `json:"journal_id" gorm:"type:uuid;not null;index"`
	Journal   JournalEntry `json:"journal" gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:CASCADE"`
	Amount    int64        `json:"amount"`
	DueDate   time.Time    `json:"due_date"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
