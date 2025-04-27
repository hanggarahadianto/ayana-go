package service

import (
	"ayana/db"
	"ayana/models"
)

func GetAssetSummaryOnly(companyID string) (int64, error) {
	var total int64

	err := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("JOIN transaction_categories ON transaction_categories.id = journal_entries.transaction_category_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ? AND transaction_categories.debit_account_type = ?", companyID, "paid", true, "Asset").
		Select("SUM(journal_lines.debit)").Scan(&total).Error

	if err != nil {
		// Jika error karena tidak ada data atau lainnya, tetap kembalikan 0 tanpa error
		return 0, nil
	}

	return total, nil
}
