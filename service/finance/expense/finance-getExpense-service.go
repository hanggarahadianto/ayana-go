package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	service "ayana/service/journalEntry"
	"ayana/utils/helper"
	"time"

	// ts "ayana/utils/helper/typesense" // Removed to fix import cycle
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ExpenseFilterParams struct {
	CompanyID       string
	Pagination      lib.Pagination
	DateFilter      lib.DateFilter
	AccountType     string // ⬅️ Tambahkan ini
	TransactionType string
	ExpenseType     string // ⬅️ Tambahkan ini
	SummaryOnly     bool
	DebitCategory   string
	CreditCategory  string
	Search          string // ⬅️ Tambahkan ini
	SortBy          string
	SortOrder       string
}

func GetExpensesFromJournalLines(params ExpenseFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {

	var (
		lines        []models.JournalLine
		total        int64
		totalExpense int64
		response     []dto.JournalEntryResponse
		now          = time.Now()
	)

	if params.Search != "" {
		results, found, err := service.SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.AccountType,
			params.DebitCategory,
			params.CreditCategory,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			&params.ExpenseType,
			params.Pagination.Page,
			params.Pagination.Limit,
		)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data pengeluaran: %w", err)
		}

		// Jika hanya summary diperlukan
		if params.SummaryOnly {
			var totalExpense int64 = 0
			for _, line := range results {
				totalExpense += int64(line.Amount)
			}
			return nil, totalExpense, int64(found), nil
		}

		var totalExpense int64 = 0
		for _, line := range results {
			totalExpense += int64(line.Amount)
		}

		return results, totalExpense, int64(found), nil
	}

	// Build base query tanpa Limit dan Offset untuk Count dan Sum
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	// Filter expense type

	baseQuery = ApplyExpenseTypeFilterToGorm(baseQuery, params.ExpenseType)

	filteredQuery, sortBy, sortOrder := helper.ApplyCommonJournalEntryFiltersToGorm(
		baseQuery,
		helper.JournalEntryFilterParams{
			DebitCategory:  params.DebitCategory,
			CreditCategory: params.CreditCategory,
			DateFilter:     params.DateFilter,
			SortBy:         params.SortBy,
			SortOrder:      params.SortOrder,
		},
		false,
	)

	// Hitung total baris
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}
	if err := filteredQuery.
		Session(&gorm.Session{}).
		Order(nil).
		Select("COALESCE(SUM(ABS(journal_lines.debit - journal_lines.credit)), 0)").
		Scan(&totalExpense).Error; err != nil {
		return nil, 0, 0, err
	}
	if params.SummaryOnly {
		return nil, totalExpense, total, nil
	}

	dataQuery := filteredQuery.
		Preload("Journal").
		Preload("Journal.TransactionCategory").
		Order(fmt.Sprintf("journal_entries.%s %s", sortBy, sortOrder)).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	response = dto.MapJournalLinesToResponse(lines, params.ExpenseType, now)

	return response, totalExpense, total, nil
}
