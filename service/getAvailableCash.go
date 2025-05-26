package service

import (
	"ayana/db"
	"ayana/models"
)

func GetAvailableCash(companyID string) (totalCashIn int64, totalCashOut int64, availableCash int64, err error) {
	// Hitung total cashin
	err = db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ?", companyID).
		Where("journal_lines.debit > 0").
		Where("journal_lines.debit_account_type = ?", "Asset").
		Where("journal_entries.transaction_type = ?", "payin").
		Where("NOT (journal_lines.credit_account_type = 'Revenue' AND journal_entries.status = 'unpaid')").
		Select("COALESCE(SUM(journal_lines.debit), 0)").
		Scan(&totalCashIn).Error
	if err != nil {
		return
	}

	// Hitung total cashout
	err = db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ?", companyID).
		Where("journal_lines.credit > 0").
		Where("journal_lines.credit_account_type = ?", "Asset").
		Where("journal_entries.is_repaid = ?", true).
		Select("COALESCE(SUM(journal_entries.amount), 0)").
		Where("journal_entries.status IN ?", []string{"paid", "done"}).
		Where("journal_entries.transaction_type = ?", "payout").
		// Where("NOT (journal_lines.credit_account_type = 'Revenue')").
		Scan(&totalCashOut).Error
	if err != nil {
		return
	}

	// totalCashOut = -totalCashOut // amount di database negatif
	availableCash = totalCashIn - totalCashOut

	return
}
