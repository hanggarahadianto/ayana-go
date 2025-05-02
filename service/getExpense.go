package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"

	"gorm.io/gorm"
)

type ExpenseFilterParams struct {
	CompanyID     string
	Pagination    helper.Pagination
	DateFilter    helper.DateFilter
	SummaryOnly   bool
	ExpenseStatus string // <-- tetap ada

}

func GetExpensesFromJournalLines(params ExpenseFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalexpense int64

	// Build base query
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ?", params.CompanyID).
		Where("journal_lines.debit_account_type = ?", "Expense").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	// Filter expense type
	switch params.ExpenseStatus {
	case "base":
		baseQuery = baseQuery.
			Where("journal_lines.credit > 0").
			Where("journal_entries.transaction_type = ?", "payout").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")

		// case "done":
		// 	baseQuery = baseQuery.
		// 		Where("journal_lines.credit > 0").
		// 		Where("journal_entries.transaction_type = ?", "payout").
		// 		Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "done", "payout")

		// case "cashout":
		// 	baseQuery = baseQuery.
		// 		Where("journal_lines.credit > 0").
		// 		Where("journal_entries.transaction_type = ?", "payout").
		// 		Where("journal_lines.credit_account_type = ?", "expense")
		// case "receivable":
		// 	baseQuery = baseQuery.
		// 		Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", false, "unpaid", "payot")
	}

	// Filter date
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed <= ?", params.DateFilter.EndDate)
	}

	// Count total
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_lines.debit - journal_lines.credit), 0)").
		Scan(&totalexpense).Error; err != nil {
		return nil, 0, 0, err
	}

	// Sum total expense (Fix query to avoid duplicate JOIN)
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_entries.amount), 0)").
		Scan(&totalexpense).Error; err != nil {
		return nil, 0, 0, err
	}

	// If summary_only = false, fetch data
	if !params.SummaryOnly {
		dataQuery := baseQuery.Session(&gorm.Session{}).
			Preload("Journal"). // Memuat relasi dengan JournalEntry
			Where("journal_entries.company_id = ?", params.CompanyID).
			Order("journal_entries.date_inputed ASC").
			Limit(params.Pagination.Limit).
			Offset(params.Pagination.Offset)

		if err := dataQuery.Find(&lines).Error; err != nil {
			return nil, 0, 0, err
		}

		var response []dto.JournalLineResponse
		for _, line := range lines {
			// Langsung menggunakan CreditAccountType yang sudah ada di JournalLine
			response = append(response, dto.JournalLineResponse{
				ID:                line.ID.String(),
				JournalEntryID:    line.JournalID.String(), // <- Use JournalID
				Invoice:           line.Journal.Invoice,    // <- Use Journal.Invoice
				Description:       line.Journal.Description,
				Partner:           line.Journal.Partner,
				Amount:            float64(line.Debit - line.Credit), // Debit - Credit
				TransactionType:   string(line.TransactionType),
				DebitAccountType:  line.DebitAccountType,
				CreditAccountType: line.CreditAccountType, // Ambil langsung dari JournalLine
				Status:            string(line.Journal.Status),
				CompanyID:         line.CompanyID.String(),
				DateInputed:       *line.Journal.DateInputed,
				DueDate:           helper.SafeDueDate(line.Journal.DueDate),
				IsRepaid:          line.Journal.IsRepaid,
				Installment:       line.Journal.Installment,
				Note:              line.Journal.Note,
			})
		}

		return response, totalexpense, total, nil
	}

	return nil, totalexpense, total, nil
}
