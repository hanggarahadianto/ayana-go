package service

import (
	"ayana/db"
	"ayana/models"
)

func GetAvailableCash(companyID string) (int64, error) {
	var totalCashIn int64
	var totalCashOut int64

	// Hitung total cashin (debit dari journal_lines, debit_account_type = Asset)
	err := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_lines.company_id = ?", companyID).
		// Where("journal_entries.status = ?", "paid").
		Where("journal_lines.debit_account_type = ?", "Asset").
		Select("COALESCE(SUM(journal_lines.debit), 0)").
		Scan(&totalCashIn).Error
	if err != nil {
		return 0, err
	}

	// Hitung total cashout (credit dari journal_lines, credit_account_type = Asset)
	err = db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_lines.company_id = ?", companyID).
		Where("journal_entries.status = ?", "paid").
		Where("journal_entries.is_repaid = ?", true).
		Where("journal_lines.credit_account_type = ?", "Asset").
		Select("COALESCE(SUM(journal_lines.credit), 0)").
		Scan(&totalCashOut).Error
	if err != nil {
		return 0, err
	}

	// Available cash = total cash in - total cash out
	return totalCashIn - totalCashOut, nil
}
