package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	service "ayana/service/journalEntry"
	"ayana/utils/helper"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type RevenueFilterParams struct {
	CompanyID       string
	Pagination      lib.Pagination
	DateFilter      lib.DateFilter
	AccountType     string // ⬅️ Tambahkan ini
	TransactionType string
	RevenueType     string
	SummaryOnly     bool
	DebitCategory   string
	CreditCategory  string
	Search          string // ⬅️ Tambahkan ini
	SortBy          string
	SortOrder       string
}

func GetRevenueFromJournalLines(params RevenueFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {

	var (
		lines        []models.JournalLine
		total        int64
		totalRevenue int64
		response     []dto.JournalEntryResponse
		now          = time.Now()
	)

	// b, _ := json.MarshalIndent(params, "", "  ")
	// log.Println("GetRevenueFromJournalLines params:", string(b))

	if params.Search != "" {
		results, _, found, err := service.SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			params.AccountType,
			params.TransactionType,
			params.RevenueType,
			params.DebitCategory,
			params.CreditCategory,
			params.Pagination.Page,
			params.Pagination.Limit,
		)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		// Jika hanya summary diperlukan
		if params.SummaryOnly {
			var totalRevenue int64 = 0
			for _, line := range results {
				totalRevenue += int64(line.Amount)
			}
			return nil, totalRevenue, int64(found), nil
		}

		var totalRevenue int64 = 0
		for _, line := range results {
			totalRevenue += int64(line.Amount)
		}

		return results, totalRevenue, int64(found), nil
	}

	paramBytes, _ := json.MarshalIndent(params, "", "  ")
	log.Println("📥 RevenueFilterParams:\n", string(paramBytes))

	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	baseQuery = ApplyRevenueTypeFilterToGorm(baseQuery, params.RevenueType)
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

	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}
	if err := filteredQuery.
		Session(&gorm.Session{}).
		Order(nil).
		Select("COALESCE(SUM(ABS(journal_lines.debit - journal_lines.credit)), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, 0, 0, err
	}
	if params.SummaryOnly {
		return nil, totalRevenue, total, nil
	}

	dataQuery := filteredQuery.
		Preload("Journal").
		Preload("Journal.TransactionCategory")

	if sortBy == "date_inputed" {
		dataQuery = dataQuery.
			Order(fmt.Sprintf("journal_entries.date_inputed %s", sortOrder)).
			Order("journal_entries.invoice DESC")
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
	response = dto.MapJournalLinesToResponse(lines, params.RevenueType, now)

	return response, totalRevenue, total, nil
}
