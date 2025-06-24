package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	service "ayana/service/journalEntry"
	"ayana/utils/helper"

	// ts "ayana/utils/helper/typesense" // Removed to fix import cycle
	"fmt"
	"log"
	"math"

	"gorm.io/gorm"
)

type ExpenseFilterParams struct {
	CompanyID      string
	Pagination     lib.Pagination
	DateFilter     lib.DateFilter
	SummaryOnly    bool
	Status         string
	DebitCategory  string
	CreditCategory string
	Search         string // ⬅️ Tambahkan ini
	SortBy         string
	SortOrder      string
}

func GetExpensesFromJournalLines(params ExpenseFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalExpense int64

	if params.Search != "" {
		results, found, err := service.SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DebitCategory,
			params.CreditCategory,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			nil,
			nil,
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
		Where("journal_entries.company_id = ?", params.CompanyID).
		Where("journal_lines.debit_account_type = ?", "Expense")

	// Filter expense type

	baseQuery = ApplyExpenseTypeFilterToGorm(baseQuery, params.Status)

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

	var response []dto.JournalEntryResponse
	for _, line := range lines {

		response = append(response, dto.JournalEntryResponse{
			ID:                      line.JournalID.String(),
			Invoice:                 line.Journal.Invoice,
			TransactionID:           line.Journal.Transaction_ID,
			TransactionCategoryID:   line.Journal.TransactionCategoryID.String(),
			TransactionCategoryName: line.Journal.TransactionCategory.Name,
			DebitCategory:           line.Journal.TransactionCategory.DebitCategory,
			CreditCategory:          line.Journal.TransactionCategory.CreditCategory,
			Description:             line.Journal.Description,
			Partner:                 line.Journal.Partner,
			Amount:                  int64(math.Abs(float64(line.Debit - line.Credit))),
			TransactionType:         string(line.TransactionType),
			DebitAccountType:        line.DebitAccountType,
			CreditAccountType:       line.CreditAccountType,
			Status:                  string(line.Journal.Status),
			CompanyID:               line.CompanyID.String(),
			DateInputed:             line.Journal.DateInputed,
			DueDate:                 lib.SafeDueDate(line.Journal.DueDate),
			IsRepaid:                line.Journal.IsRepaid,
			Installment:             line.Journal.Installment,
			Note:                    line.Journal.Note,
		})
	}

	return response, totalExpense, total, nil
}
