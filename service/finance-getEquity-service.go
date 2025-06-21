package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"fmt"
	"log"
	"math"

	"gorm.io/gorm"
)

type EquityFilterParams struct {
	CompanyID       string
	Pagination      helper.Pagination
	DateFilter      helper.DateFilter
	SummaryOnly     bool
	EquityType      string
	TransactionType string
	DebitCategory   string
	CreditCategory  string
	Search          string // ⬅️ Tambahkan ini
}

func GetEquityFromJournalLines(params EquityFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalEquity int64

	if params.Search != "" {
		results, found, err := SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DebitCategory,
			params.CreditCategory,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
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

	switch params.EquityType {
	case "setor":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Equity")

		// Where("NOT (journal_lines.credit_account_type = 'Revenue' AND journal_entries.status = 'unpaid')")

	case "tarik":
		baseQuery = baseQuery.
			Where("journal_lines.credit > 0").
			Where("journal_entries.transaction_type = ?", "payout").
			Where("journal_lines.debit_account_type = ?", "Equity").
			Where("journal_lines.credit_account_type = ?", "Asset")

	}

	if params.DebitCategory != "" {
		baseQuery = baseQuery.Where("LOWER(transaction_categories.debit_category) = LOWER(?)", params.DebitCategory)

	}
	if params.CreditCategory != "" {
		baseQuery = baseQuery.Where("LOWER(transaction_categories.credit_category) = LOWER(?)", params.CreditCategory)

	}

	// Filter date
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.
			Where("journal_entries.date_inputed BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.
			Where("journal_entries.date_inputed >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.
			Where("journal_entries.date_inputed <= ?", params.DateFilter.EndDate)
	}

	// Hitung total baris
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Hitung total asset (menggunakan debit - credit)
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_lines.debit - journal_lines.credit), 0)").
		Scan(&totalEquity).Error; err != nil {
		return nil, 0, 0, err
	}

	// Query untuk mengambil data dengan pagination
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Preload("Journal").
		Preload("Journal.TransactionCategory"). // ✅ Tambahkan ini
		Order("journal_entries.date_inputed ASC").
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
			DueDate:                 helper.SafeDueDate(line.Journal.DueDate),
			IsRepaid:                line.Journal.IsRepaid,
			Installment:             line.Journal.Installment,
			Note:                    line.Journal.Note,
		})
	}

	return response, totalEquity, total, nil
}
