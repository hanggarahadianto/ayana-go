package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"

	"gorm.io/gorm"
)

type CashFilterParams struct {
	CompanyID  string
	Status     string
	Pagination helper.Pagination
	DateFilter helper.DateFilter
}

// GetCash ambil list cash entry
func GetCash(params CashFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var entries []models.JournalEntry

	// 1. Base Query
	baseQuery := db.DB.Model(&models.JournalEntry{}).
		Where("company_id = ? AND transaction_type = ?", params.CompanyID, "payin").
		Where("date_inputed BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)

	// 2. Hitung total record
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 3. Ambil list data + order by due_date
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Order("due_date ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&entries).Error; err != nil {
		return nil, 0, 0, err
	}

	// 4. Mapping ke DTO
	var responseData []dto.JournalEntryResponse
	for _, entry := range entries {
		responseData = append(responseData, dto.JournalEntryResponse{
			ID:                    entry.ID.String(),
			Invoice:               entry.Invoice,
			Description:           entry.Description,
			TransactionCategoryID: entry.TransactionCategoryID.String(),
			Amount:                float64(entry.Amount),
			Partner:               entry.Partner,
			TransactionType:       string(entry.TransactionType),
			Status:                string(entry.Status),
			CompanyID:             entry.CompanyID.String(),
			DateInputed:           *entry.DateInputed,
			DueDate:               *entry.DueDate,
			IsRepaid:              entry.IsRepaid,
			Installment:           entry.Installment,
			Note:                  entry.Note,
		})
	}

	return responseData, total, int64(len(responseData)), nil
}
