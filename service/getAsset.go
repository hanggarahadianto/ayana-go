package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"

	"gorm.io/gorm"
)

type AssetFilterParams struct {
	CompanyID       string
	Pagination      helper.Pagination
	DateFilter      helper.DateFilter
	SummaryOnly     bool
	AssetType       string // <-- tetap ada
	TransactionType string // <-- tetap ada
}

func GetAssetsFromJournalLines(params AssetFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalAsset int64

	// Build base query
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ?", params.CompanyID).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	// Filter asset type
	switch params.AssetType {
	case "cashin":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset")

	case "fixed_asset":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")

	case "cashout":
		baseQuery = baseQuery.
			Where("journal_lines.credit > 0").
			Where("journal_lines.credit_account_type = ?", "Asset").
			// Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")
			Where("journal_entries.is_repaid = ? AND journal_entries.status IN (?, ?) AND journal_entries.transaction_type = ?", true, "paid", "done", "payout")
	case "receivable":
		baseQuery = baseQuery.
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", false, "unpaid", "payot")
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
		Scan(&totalAsset).Error; err != nil {
		return nil, 0, 0, err
	}

	// Sum total asset (Fix query to avoid duplicate JOIN)
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_entries.amount), 0)").
		Scan(&totalAsset).Error; err != nil {
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
				Transaction_ID:    line.Journal.Transaction_ID,
				Partner:           line.Journal.Partner,
				Invoice:           line.Journal.Invoice, // <- Use Journal.Invoice
				Description:       line.Journal.Description,
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

		return response, totalAsset, total, nil
	}

	return nil, totalAsset, total, nil
}
