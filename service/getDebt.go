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

type DebtFilterParams struct {
	CompanyID   string
	Pagination  helper.Pagination
	DateFilter  helper.DateFilter
	SummaryOnly bool
	DebtStatus  string
	Category    string
	Search      string // ⬅️ Tambahkan ini
}

func GetDebtsFromJournalLines(params DebtFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalDebt int64

	if params.Search != "" {
		results, found, err := SearchJournalLines(params.Search, params.CompanyID, params.Category, params.Pagination.Page, params.Pagination.Limit)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		// Jika hanya summary diperlukan
		if params.SummaryOnly {
			var totalDebt int64 = 0
			for _, line := range results {
				totalDebt += int64(line.Amount)
			}
			return nil, totalDebt, int64(found), nil
		}

		var totalDebt int64 = 0
		for _, line := range results {
			totalDebt += int64(line.Amount)
		}

		return results, totalDebt, int64(found), nil
	}

	// Build base query
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	switch params.DebtStatus {
	case "going":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ?", false, "unpaid").
			Where("LOWER(journal_lines.credit_account_type) = ?", "liability").
			Where("LOWER(journal_lines.debit_account_type) != ?", "revenue")
	case "done":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ?", true, "done").
			Where("LOWER(journal_lines.debit_account_type) = ?", "liability").
			Where("LOWER(journal_lines.credit_account_type) = ?", "asset")
	}

	if params.Category != "" {
		baseQuery = baseQuery.Where("transaction_categories.category ILIKE ?", "%"+params.Category+"%")
	}

	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed <= ?", params.DateFilter.EndDate)
	}

	// Total count
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Total debt
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_entries.amount), 0)").
		Scan(&totalDebt).Error; err != nil {
		return nil, 0, 0, err
	}

	if params.SummaryOnly {
		return nil, totalDebt, total, nil
	}

	// Data with pagination
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Preload("Journal").
		Preload("Journal.TransactionCategory").
		Order("journal_entries.date_inputed ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	var response []dto.JournalLineResponse
	for _, line := range lines {
		var paymentDateStatus string

		// Hitung payment status jika DebtStatus adalah "done"
		if params.DebtStatus == "done" && line.Journal.DateInputed != nil && line.Journal.DueDate != nil {
			due := *line.Journal.DueDate
			input := *line.Journal.DateInputed

			// Lewati jika due date kosong (zero time)
			if !due.IsZero() {
				diff := due.Sub(input).Hours() / 24
				if diff >= 0 {
					paymentDateStatus = fmt.Sprintf("Dibayar Tepat Waktu %.0f Hari Sebelum Jatuh Tempo", diff)
				} else {
					paymentDateStatus = fmt.Sprintf("Dibayar Terlambat %.0f Hari Setelah Jatuh Tempo", -diff)
				}
			}
		}

		response = append(response, dto.JournalLineResponse{
			ID:                      line.JournalID.String(),
			Transaction_ID:          line.Journal.Transaction_ID,
			TransactionCategoryID:   line.Journal.TransactionCategoryID.String(),
			TransactionCategoryName: line.Journal.TransactionCategory.Name,
			Category:                line.Journal.TransactionCategory.Category,
			Invoice:                 line.Journal.Invoice,
			Description:             line.Journal.Description,
			Partner:                 line.Journal.Partner,
			Amount:                  math.Abs(float64(line.Debit - line.Credit)),
			TransactionType:         string(line.TransactionType),
			DebitAccountType:        line.DebitAccountType,
			CreditAccountType:       line.CreditAccountType,
			Status:                  string(line.Journal.Status),
			CompanyID:               line.CompanyID.String(),
			DateInputed:             *line.Journal.DateInputed,
			DueDate:                 helper.SafeDueDate(line.Journal.DueDate),
			IsRepaid:                line.Journal.IsRepaid,
			Installment:             line.Journal.Installment,
			Note:                    line.Journal.Note,
			PaymentDateStatus:       paymentDateStatus, // <-- Tambahan field baru
		})
	}

	return response, totalDebt, total, nil
}
