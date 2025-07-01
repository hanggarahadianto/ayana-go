package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	service "ayana/service/journalEntry"
	"ayana/utils/helper"

	"fmt"
	"time"

	"gorm.io/gorm"
)

type AssetFilterParams struct {
	CompanyID       string
	Pagination      lib.Pagination
	DateFilter      lib.DateFilter
	AccountType     string // ⬅️ Tambahkan ini
	TransactionType string
	AssetType       string
	SummaryOnly     bool
	DebitCategory   string
	CreditCategory  string
	Search          string
	SortBy          string
	SortOrder       string
}

func GetAssetsFromJournalLines(params AssetFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var (
		lines      []models.JournalLine
		total      int64
		totalAsset int64
		response   []dto.JournalEntryResponse
		now        = time.Now()
	)

	if params.Search != "" {
		results, _, found, err := service.SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			params.AccountType,
			params.TransactionType,
			params.AssetType,
			params.DebitCategory,
			params.CreditCategory,
			params.Pagination.Page,
			params.Pagination.Limit,
		)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		// ✅ Hitung totalAsset langsung dari hasil search
		var totalAsset int64
		for _, r := range results {
			totalAsset += int64(r.Amount)
		}

		results = helper.EnrichJournalEntryResponses(results, params.AssetType, now)
		return results, totalAsset, int64(found), nil
	}

	// ✅ Build base query
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	// ✅ Apply filters (tanpa sort dulu)
	baseQuery = ApplyAssetTypeFilterToGorm(baseQuery, params.AssetType)
	filteredQuery, sortBy, sortOrder := helper.ApplyCommonJournalEntryFiltersToGorm(
		baseQuery,
		helper.JournalEntryFilterParams{
			DebitCategory:  params.DebitCategory,
			CreditCategory: params.CreditCategory,
			DateFilter:     params.DateFilter,
			SortBy:         params.SortBy,
			SortOrder:      params.SortOrder,
		},
		false, // no sort yet
	)

	// ✅ Count total
	if err := filteredQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// ✅ Hitung total asset (hindari ORDER BY)
	if err := filteredQuery.
		Session(&gorm.Session{}).
		Order(nil).
		Select("COALESCE(SUM(ABS(journal_lines.debit - journal_lines.credit)), 0)").
		Scan(&totalAsset).Error; err != nil {
		return nil, 0, 0, err
	}

	if params.SummaryOnly {
		return nil, totalAsset, total, nil
	}

	// ✅ Apply pagination dan sorting
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
	response = dto.MapJournalLinesToResponse(lines, params.AssetType, now)

	return response, totalAsset, total, nil
}
