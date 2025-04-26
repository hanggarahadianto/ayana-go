package service

import (
	"ayana/db"
	"ayana/models"
)

func GetExpenseSummaryOnly(companyID string) (int64, error) {
	var totalExpense int64

	err := db.DB.Model(&models.JournalLine{}).
		Select("COALESCE(SUM(journal_lines.debit), 0)").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("JOIN transaction_categories ON transaction_categories.id = journal_entries.transaction_category_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ? AND transaction_categories.debit_account_type = ? AND journal_lines.debit > 0",
			companyID, "paid", true, "Expense").
		Scan(&totalExpense).Error

	if err != nil {
		return 0, err
	}

	return totalExpense, nil
}
