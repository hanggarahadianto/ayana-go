package service

import (
	"ayana/db"
	"ayana/models"
)

func GetExpenseSummaryOnly(companyID string) (int64, error) {
	var totalExpense int64

	err := db.DB.Model(&models.JournalLine{}).
		Select("COALESCE(SUM(journal_lines.credit), 0)").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("journal_entries.company_id = ? AND journal_entries.transaction_type = ? AND journal_entries.status = ? AND journal_lines.credit > 0",
			companyID, "payout", "paid").
		Scan(&totalExpense).Error

	if err != nil {
		return 0, err // hanya jika query gagal, bukan karena NULL result
	}

	return totalExpense, nil
}
