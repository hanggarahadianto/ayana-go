package service

import "gorm.io/gorm"

func ApplyEquityTypeFilterToGorm(query *gorm.DB, EquityType string) *gorm.DB {

	switch EquityType {
	case "setor":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Equity")

	case "tarik":
		return query.
			Where("journal_lines.credit > 0").
			Where("journal_entries.transaction_type = ?", "payout").
			Where("journal_lines.debit_account_type = ?", "Equity").
			Where("journal_lines.credit_account_type = ?", "Asset")

	default:
		return query
	}
}
