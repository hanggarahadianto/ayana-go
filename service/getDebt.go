package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"

	"gorm.io/gorm"
)

type DebtFilterParams struct {
	CompanyID string

	Pagination  helper.Pagination
	DateFilter  helper.DateFilter
	SummaryOnly bool // Menambahkan SummaryOnly
}

// Mengambil list + summary debt
func GetOutstandingDebts(params DebtFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var entries []models.JournalEntry
	var total int64
	var totalDebt int64

	// SubQuery untuk filter journal_id
	subQuery := db.DB.
		Model(&models.JournalLine{}).
		Select("journal_id").
		Where("credit > 0 AND company_id = ?", params.CompanyID)

	baseQuery := db.DB.Model(&models.JournalEntry{}).
		Where("id IN (?) AND status = ? AND is_repaid = false", subQuery, "unpaid")

	// Tambahkan filter tanggal kalau ada
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("due_date >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("due_date <= ?", params.DateFilter.EndDate)
	}

	// 1. Hitung total jumlah data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 2. Hitung total hutang (tanpa ORDER BY, biar ga error)
	if err := baseQuery.Session(&gorm.Session{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalDebt).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jika SummaryOnly = true, kembalikan hanya total dan totalDebt, tanpa mengambil data list
	if params.SummaryOnly {
		return nil, totalDebt, total, nil
	}

	// 3. Ambil list data + order by due_date (jika SummaryOnly = false)
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Order("due_date ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&entries).Error; err != nil {
		return nil, 0, 0, err
	}

	// Mapping ke DTO
	var responseData []dto.JournalEntryResponse
	for _, entry := range entries {
		responseData = append(responseData, dto.JournalEntryResponse{
			ID:                    entry.ID.String(),
			TransactionID:         entry.Transaction_ID,
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

	return responseData, totalDebt, total, nil
}
