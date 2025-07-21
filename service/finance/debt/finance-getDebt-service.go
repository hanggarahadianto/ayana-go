package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	service "ayana/service/journalEntry"
	"ayana/utils/helper"

	"ayana/models"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type DebtFilterParams struct {
	CompanyID       string
	Pagination      lib.Pagination
	DateFilter      lib.DateFilter
	AccountType     string // ⬅️ Tambahkan ini
	TransactionType string
	DebtType        string // ⬅️ Tambahkan ini
	DebitCategory   string
	CreditCategory  string
	SummaryOnly     bool
	Search          string
	SortBy          string
	SortOrder       string
}

func GetDebtsFromJournalLines(params DebtFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var (
		lines     []models.JournalLine
		total     int64
		totalDebt int64
		response  []dto.JournalEntryResponse
		now       = time.Now()
	)

	// 🔍 Handle Search via Typesense
	if params.Search != "" {
		results, _, found, err := service.SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			params.AccountType,
			params.TransactionType,
			params.DebtType,
			params.DebitCategory,
			params.CreditCategory,
			params.Pagination.Page,
			params.Pagination.Limit,
		)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data hutang: %w", err)
		}

		for i, line := range results {
			note, color := lib.HitungPaymentNote(line.DueDate, line.RepaymentDate, params.DebtType, now)
			results[i].PaymentNote = note
			results[i].PaymentNoteColor = color
			totalDebt += int64(line.Amount)
		}

		if params.SummaryOnly {
			return nil, totalDebt, int64(found), nil
		}

		return results, totalDebt, int64(found), nil
	}

	// paramBytes, _ := json.MarshalIndent(params, "", "  ")
	// log.Println("📥 DebtFilterParams:\n", string(paramBytes))

	// 🏗️ Build base query
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	// 📦 Filter khusus hutang
	baseQuery = ApplyDebtTypeFilterToGorm(baseQuery, params.DebtType)

	// 🎛️ Filter umum (tanpa sorting dulu)
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

	// 🔢 Hitung total data
	if err := filteredQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 💰 Hitung total nilai hutang (tanpa ORDER BY)
	if err := filteredQuery.
		Session(&gorm.Session{}).
		Order(nil).
		Select("COALESCE(SUM(ABS(journal_lines.debit - journal_lines.credit)), 0)").
		Scan(&totalDebt).Error; err != nil {
		return nil, 0, 0, err
	}

	if params.SummaryOnly {
		return nil, totalDebt, total, nil
	}

	dataQuery := filteredQuery.
		Preload("Journal").
		Preload("Journal.TransactionCategory")

	if sortBy == "date_inputed" {
		dataQuery = dataQuery.
			Order(fmt.Sprintf("journal_entries.date_inputed %s", sortOrder)). // Primary sort
			Order("journal_entries.invoice DESC")                             // Secondary sort, fixed DESC
	} else {
		dataQuery = dataQuery.
			Order(fmt.Sprintf("journal_entries.%s %s", sortBy, sortOrder))
	}

	dataQuery = dataQuery.
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	// 🧾 Mapping response
	response = dto.MapJournalLinesToResponse(lines, params.DebtType, now)

	return response, totalDebt, total, nil
}
