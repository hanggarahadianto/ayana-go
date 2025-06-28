package service

import (
	"fmt"

	"gorm.io/gorm"
)

func ApplyExpenseTypeFilterToGorm(query *gorm.DB, expenseType string) *gorm.DB {
	fmt.Println("🔥 ApplyExpenseTypeFilterToGorm triggered with expense type:", expenseType)

	switch expenseType {
	case "base":
		return query.
			Where("journal_lines.debit > 0").
			Where("LOWER(journal_lines.debit_account_type) = ?", "expense").
			Where("LOWER(journal_lines.credit_account_type) = ?", "asset").
			Where("journal_entries.is_repaid = ?", true).
			Where("LOWER(journal_entries.status) IN ?", []string{"paid", "done"}).
			Where("LOWER(journal_entries.transaction_type) = ?", "payout")
	default:
		return query
	}
}
