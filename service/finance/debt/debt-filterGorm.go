package service

import "gorm.io/gorm"

func ApplyDebtTypeFilterToGorm(query *gorm.DB, debtType string) *gorm.DB {

	switch debtType {
	case "going":
		return query.
			Where("journal_lines.credit > 0").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ?", false, "unpaid").
			Where("LOWER(journal_lines.credit_account_type) = ?", "liability").
			Where("LOWER(journal_lines.debit_account_type) != ?", "revenue")

	case "done":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ?", true, "done").
			Where("LOWER(journal_lines.debit_account_type) = ?", "liability").
			Where("LOWER(journal_lines.credit_account_type) = ?", "asset")

	default:
		return query
	}
}
