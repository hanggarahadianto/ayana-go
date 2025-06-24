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
	"math"
	"time"

	"gorm.io/gorm"
)

type AssetFilterParams struct {
	CompanyID       string
	Pagination      lib.Pagination
	DateFilter      lib.DateFilter
	SummaryOnly     bool
	AssetType       string
	Status          string
	TransactionType string
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

	// ✅ Handle Typesense Search
	if params.Search != "" {
		results, found, err := service.SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DebitCategory,
			params.CreditCategory,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			&params.AssetType,
			&params.Status,
			params.Pagination.Page,
			params.Pagination.Limit,
		)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		for i, line := range results {
			note, color := lib.HitungPaymentNote(params.AssetType, line.DueDate, line.RepaymentDate, now)
			results[i].PaymentNote = note
			results[i].PaymentNoteColor = color
			totalAsset += int64(line.Amount)
		}

		if params.SummaryOnly {
			return nil, totalAsset, int64(found), nil
		}

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

	// ✅ Map ke DTO
	for _, line := range lines {
		note, color := lib.HitungPaymentNote(params.AssetType, line.Journal.DueDate, line.Journal.RepaymentDate, now)

		entry := dto.JournalEntryResponse{
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
			RepaymentDate:           line.Journal.RepaymentDate,
			IsRepaid:                line.Journal.IsRepaid,
			Installment:             line.Journal.Installment,
			Note:                    line.Journal.Note,
		}

		if params.AssetType == "receivable" {
			entry.PaymentNote = note
			entry.PaymentNoteColor = color
		}

		response = append(response, entry)
	}

	return response, totalAsset, total, nil
}
