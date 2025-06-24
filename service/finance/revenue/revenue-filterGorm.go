package service

import "gorm.io/gorm"

func ApplyRevenueTypeFilterToGorm(query *gorm.DB, Status string) *gorm.DB {

	switch Status {
	case "realized":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Revenue").
			Where("journal_entries.status IN ?", []string{"done", "paid"})

	case "unrealized":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Revenue").
			Where("NOT (journal_entries.status IN ?)", []string{"done", "paid"})

	default:
		return query
	}
}
