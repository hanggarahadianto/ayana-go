package service

import (
	"fmt"

	"gorm.io/gorm"
)

func ApplyExpenseTypeFilterToGorm(query *gorm.DB, Status string) *gorm.DB {
	fmt.Println("🔥 ApplyExpenseTypeFilterToGorm triggered with status:", Status)

	switch Status {
	case "base":
		return query.
			Where("journal_lines.debit > 0").
			Where("LOWER(journal_lines.credit_account_type) = ?", "asset").
			Where("journal_entries.is_repaid = ? AND journal_entries.status IN ? AND journal_entries.transaction_type = ?", true, []string{"paid", "done"}, "payout")

	default:
		return query
	}
}
