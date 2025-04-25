package dto

import "github.com/google/uuid"

type AccountResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        int16     `json:"code"` // <- ubah dari int ke int16
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CompanyID   uuid.UUID `json:"company_id"`
}
