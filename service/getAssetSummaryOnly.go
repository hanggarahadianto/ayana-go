package service

import (
	"ayana/db"
	"ayana/models"
	"errors"
)

func GetAssetSummaryOnly(companyID string) (int64, error) {
	var total int64

	// Menghitung total asset (penjumlahan amount seluruh data tanpa paginasi)
	err := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("JOIN transaction_categories ON transaction_categories.id = journal_entries.transaction_category_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ? AND transaction_categories.debit_account_type = ?", companyID, "paid", true, "Asset").
		Select("SUM(journal_lines.debit)").Scan(&total).Error

	if err != nil {
		return 0, errors.New("failed to calculate summary total")
	}

	return total, nil
}
