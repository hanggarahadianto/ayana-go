package service

import "gorm.io/gorm"

func ApplyExpenseTypeFilterToGorm(query *gorm.DB, Status string) *gorm.DB {

	switch Status {
	case "base":
		return query.
			Where("journal_lines.debit > 0"). // Expense biasanya di sisi debit
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")

	default:
		return query
	}
}
