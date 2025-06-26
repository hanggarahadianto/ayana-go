package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	service "ayana/service/journalEntry"
	"ayana/utils/helper"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type EquityFilterParams struct {
	CompanyID       string
	Pagination      lib.Pagination
	DateFilter      lib.DateFilter
	AccountType     string // ⬅️ Tambahkan ini
	SummaryOnly     bool
	EquityType      string
	TransactionType string
	DebitCategory   string
	CreditCategory  string
	Search          string // ⬅️ Tambahkan ini
	SortBy          string
	SortOrder       string
}

func GetEquityFromJournalLines(params EquityFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var (
		lines       []models.JournalLine
		total       int64
		totalEquity int64
		response    []dto.JournalEntryResponse
		now         = time.Now()
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
			nil,
			params.Pagination.Page,
			params.Pagination.Limit,
		)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		// Jika hanya summary diperlukan
		if params.SummaryOnly {
			var totalEquity int64 = 0
			for _, line := range results {
				totalEquity += int64(line.Amount)
			}
			return nil, totalEquity, int64(found), nil
		}

		var totalEquity int64 = 0
		for _, line := range results {
			totalEquity += int64(line.Amount)
		}

		return results, totalEquity, int64(found), nil
	}

	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	baseQuery = ApplyEquityTypeFilterToGorm(baseQuery, params.EquityType)

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

	if err := filteredQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Hitung total asset (menggunakan debit - credit)
	if err := filteredQuery.
		Session(&gorm.Session{}).
		Order(nil).
		Select("COALESCE(SUM(ABS(journal_lines.debit - journal_lines.credit)), 0)").
		Scan(&totalEquity).Error; err != nil {
		return nil, 0, 0, err
	}
	if params.SummaryOnly {
		return nil, totalEquity, total, nil
	}

	// Query untuk mengambil data dengan pagination
	dataQuery := filteredQuery.
		Preload("Journal").
		Preload("Journal.TransactionCategory").
		Order(fmt.Sprintf("journal_entries.%s %s", sortBy, sortOrder)).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	// 🧾 Mapping response
	response = dto.MapJournalLinesToResponse(lines, params.EquityType, now)

	return response, totalEquity, total, nil
}
