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
	AssetType       string
	TransactionType string
	Category        string
}

func GetAssetsFromJournalLines(params AssetFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalAsset int64

	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID)

	switch params.AssetType {
	case "cashin":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("NOT (journal_lines.credit_account_type = 'Revenue' AND journal_entries.status = 'unpaid')")

	case "fixed_asset":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")

	case "cashout":
		baseQuery = baseQuery.
			Where("journal_lines.credit > 0").
			Where("journal_lines.credit_account_type = ?", "Asset").
			Where("journal_entries.is_repaid = ?", true).
			Where("journal_entries.status IN ?", []string{"paid", "done"}).
			Where("journal_entries.transaction_type = ?", "payout")

	case "receivable":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Revenue").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", false, "unpaid", "payin")
		// Where("LOWER(transaction_category.category) ILIKE ?", "%piutang%") // Hanya untuk 'receivable'

	}

	if params.Category != "" {
		baseQuery = baseQuery.Where("transaction_categories.category ILIKE ?", "%"+params.Category+"%")
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

	// Hitung total asset (menggunakan debit - credit)
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_lines.debit - journal_lines.credit), 0)").
		Scan(&totalAsset).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jika SummaryOnly = true, kembalikan hanya totalAsset dan total
	if params.SummaryOnly {
		return nil, totalAsset, total, nil
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

	var response []dto.JournalLineResponse
	for _, line := range lines {
		response = append(response, dto.JournalLineResponse{
			ID:                line.ID.String(),
			JournalEntryID:    line.JournalID.String(),
			Transaction_ID:    line.Journal.Transaction_ID,
			Category:          line.Journal.TransactionCategory.Category,
			Partner:           line.Journal.Partner,
			Invoice:           line.Journal.Invoice,
			Description:       line.Journal.Description,
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

	return response, totalAsset, total, nil
}
