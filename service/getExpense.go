package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"time"

	"gorm.io/gorm"
)

// Mengambil list dan summary pengeluaran
// ExpenseFilterParams adalah parameter untuk filter pengeluaran
type ExpenseFilterParams struct {
	CompanyID   string
	Pagination  helper.Pagination
	DateFilter  helper.DateFilter
	SummaryOnly bool // Menambahkan parameter SummaryOnly
}

// Mengambil list dan summary pengeluaran
func GetExpenses(params ExpenseFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var entries []models.JournalEntry
	var total int64
	var totalExpense int64

	// SubQuery untuk filter journal_id
	subQuery := db.DB.
		Model(&models.JournalLine{}).
		Select("journal_id").
		Where("credit > 0 AND company_id = ?", params.CompanyID).
		Where("debit_account_type = ?", "Expense")

	baseQuery := db.DB.Model(&models.JournalEntry{}).
		Where("id IN (?) AND status = ? AND transaction_type = ? AND is_repaid = ?", subQuery, "paid", "payout", true)

	// Tambahkan filter tanggal kalau ada
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("date_inputed <= ?", params.DateFilter.EndDate)
	}

	// 1. Hitung total jumlah data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 2. Hitung total pengeluaran (tanpa ORDER BY, biar ga error)
	if err := baseQuery.Session(&gorm.Session{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jika summaryOnly = true, hanya kembalikan total dan totalExpense tanpa mengambil data list
	if params.SummaryOnly {
		return nil, totalExpense, total, nil
	}

	// 3. Ambil list data + order by date_inputed (jika summaryOnly = false)
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Order("date_inputed ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&entries).Error; err != nil {
		return nil, 0, 0, err
	}

	// Mapping ke DTO
	var responseData []dto.JournalEntryResponse
	for _, entry := range entries {
		var dueDate time.Time
		if entry.DueDate != nil {
			dueDate = *entry.DueDate // Menggunakan nilai sebenarnya
		} else {
			dueDate = time.Time{} // Nilai default untuk time.Time, yaitu zero value
		}

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
			DueDate:               dueDate, // Sekarang menggunakan tipe time.Time
			IsRepaid:              entry.IsRepaid,
			Installment:           entry.Installment,
			Note:                  entry.Note,
		})
	}

	return responseData, totalExpense, total, nil
}
