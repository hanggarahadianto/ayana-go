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
	ExpenseStatus string
}

func GetExpensesFromJournalLines(params ExpenseFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalExpense int64

	// Build base query tanpa Limit dan Offset untuk Count dan Sum
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ?", params.CompanyID).
		Where("journal_lines.debit_account_type = ?", "Expense")

	// Filter expense type
	switch params.ExpenseStatus {
	case "base":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0"). // Expense biasanya di sisi debit
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")
	}

	// Filter date
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed <= ?", params.DateFilter.EndDate)
	}

	// Hitung total baris
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Hitung total expense (menggunakan debit - credit)
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_lines.debit - journal_lines.credit), 0)").
		Scan(&totalExpense).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jika SummaryOnly = true, kembalikan hanya totalExpense dan total
	if params.SummaryOnly {
		return nil, totalExpense, total, nil
	}

	// Query untuk mengambil data dengan pagination
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Preload("Journal").
		Order("journal_entries.date_inputed ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	var response []dto.JournalLineResponse
	for _, line := range lines {
		response = append(response, dto.JournalLineResponse{
			ID:                line.JournalID.String(),
			Invoice:           line.Journal.Invoice,
			Transaction_ID:    line.Journal.Transaction_ID,
			Description:       line.Journal.Description,
			Partner:           line.Journal.Partner,
			Amount:            -float64(line.Debit - line.Credit),
			TransactionType:   string(line.TransactionType),
			DebitAccountType:  line.DebitAccountType,
			CreditAccountType: line.CreditAccountType,
			Status:            string(line.Journal.Status),
			CompanyID:         line.CompanyID.String(),
			DateInputed:       *line.Journal.DateInputed,
			DueDate:           helper.SafeDueDate(line.Journal.DueDate),
			IsRepaid:          line.Journal.IsRepaid,
			Installment:       line.Journal.Installment,
			Note:              line.Journal.Note,
		})
	}

	return response, totalExpense, total, nil
}
