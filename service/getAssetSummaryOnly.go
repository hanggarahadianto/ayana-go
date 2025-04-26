package service

import (
	"ayana/db"
	"ayana/models"
)

func GetAssetSummaryOnly(companyID string) (int64, error) {
	var total int64
	err := db.DB.
		Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("JOIN transaction_categories ON transaction_categories.id = journal_entries.transaction_category_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ? AND transaction_categories.debit_account_type = ?", companyID, "paid", true, "Asset").
		Scan(&total).Error
	return total, err
}
