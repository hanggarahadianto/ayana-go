package service

import (
	"ayana/db"
	"ayana/models"
)

func GetAssetSummaryOnly(companyID string) (int64, error) {
	var total int64
	err := db.DB.
		Model(&models.JournalLine{}).
		Select("COALESCE(SUM(journal_lines.debit), 0)").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ?", companyID, "paid", true).
		Scan(&total).Error
	return total, err
}
