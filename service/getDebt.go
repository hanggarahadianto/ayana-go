package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"

	"gorm.io/gorm"
)

type DebtFilterParams struct {
	CompanyID   string
	Pagination  helper.Pagination
	DateFilter  helper.DateFilter
	SummaryOnly bool
	DebtStatus  string // <-- tetap ada

}

func GetDebtsFromJournalLines(params DebtFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalDebt int64

	// Build base query tanpa Limit dan Offset untuk Count dan Sum
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	// Filter Debt type
	switch params.DebtStatus {
	case "going":
		baseQuery = baseQuery.
			Where("journal_lines.credit > 0").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ?", false, "unpaid").
			// Filter untuk memastikan bahwa hutang terdaftar di akun Liability, bukan Revenue
			Where("LOWER(journal_lines.credit_account_type) = ?", "liability").
			Where("LOWER(journal_lines.debit_account_type) != ?", "revenue")

	case "done":
		baseQuery = baseQuery.
			Where("journal_lines.credit > 0").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ?", true, "done").
			// Filter untuk memastikan bahwa hutang terdaftar di akun Liability, bukan Revenue
			Where("LOWER(journal_lines.credit_account_type) = ?", "liability").
			Where("LOWER(journal_lines.debit_account_type) != ?", "revenue")
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

	// Hitung total debt
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_entries.amount), 0)").
		Scan(&totalDebt).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jika SummaryOnly = true, kembalikan hanya total dan totalDebt
	if params.SummaryOnly {
		return nil, totalDebt, total, nil
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
			ID:                line.ID.String(),
			JournalEntryID:    line.JournalID.String(),
			Transaction_ID:    line.Journal.Transaction_ID,
			Invoice:           line.Journal.Invoice,
			Description:       line.Journal.Description,
			Partner:           line.Journal.Partner,
			Amount:            float64(line.Debit - line.Credit),
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

	return response, totalDebt, total, nil
}
