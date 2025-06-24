package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	"fmt"
	"log"
	"math"
	"time"

	"gorm.io/gorm"
)

type DebtFilterParams struct {
	CompanyID      string
	Pagination     lib.Pagination
	DateFilter     lib.DateFilter
	SummaryOnly    bool
	Status         string
	DebtStatus     string
	DebitCategory  string
	CreditCategory string
	Search         string // ⬅️ Tambahkan ini
}

func GetDebtsFromJournalLines(params DebtFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalDebt int64

	if params.Search != "" {
		results, found, err := SearchJournalLines(
			params.Search,
			params.CompanyID,
			params.DebitCategory,
			params.CreditCategory,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			nil,
			&params.Status,
			params.Pagination.Page,
			params.Pagination.Limit,
		)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		now := time.Now()

		// Tambahkan logic paymentNote ke hasil dari typesense
		for i, line := range results {
			note, color := lib.HitungPaymentNote(params.DebtStatus, line.DueDate, line.RepaymentDate, now)
			results[i].PaymentNote = note
			results[i].PaymentNoteColor = color
		}

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
			Where("journal_lines.credit > 0").
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

	if params.DebitCategory != "" {
		baseQuery = baseQuery.Where("LOWER(transaction_categories.debit_category) = LOWER(?)", params.DebitCategory)

	}
	if params.CreditCategory != "" {
		baseQuery = baseQuery.Where("LOWER(transaction_categories.credit_category) = LOWER(?)", params.CreditCategory)

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
		Order("journal_entries.due_date ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	now := time.Now()
	var response []dto.JournalEntryResponse

	for _, line := range lines {

		note, color := lib.HitungPaymentNote(params.DebtStatus, line.Journal.DueDate, line.Journal.RepaymentDate, now)

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
			RepaymentDate:           line.Journal.RepaymentDate,
			IsRepaid:                line.Journal.IsRepaid,
			Installment:             line.Journal.Installment,
			Note:                    line.Journal.Note,
			PaymentNote:             note,
			PaymentNoteColor:        color,
		})
	}

	return response, totalDebt, total, nil
}
