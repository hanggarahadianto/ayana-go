package service

import "gorm.io/gorm"

func ApplyAssetTypeFilterToGorm(query *gorm.DB, assetType string) *gorm.DB {
	switch assetType {
	case "cashin":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("NOT (journal_lines.credit_account_type = 'Revenue' AND journal_entries.status = 'unpaid')")

	case "cashout":
		return query.
			Where("journal_lines.credit > 0").
			Where("journal_lines.credit_account_type = ?", "Asset").
			Where("journal_entries.transaction_type = ?", "payout").
			Where("journal_entries.is_repaid = ?", true).
			Where("journal_entries.status IN ?", []string{"paid", "done"})

	case "fixed_asset":
		return query.
			Where("journal_lines.debit > 0").
			Where("transaction_categories.debit_account_type = ?", "Asset").
			Where("transaction_categories.debit_category = ?", "Aset Tetap").
			Where("transaction_categories.credit_account_type IN ?", []string{"Asset", "Liability"}).
			Where("transaction_categories.transaction_type = ?", "payout").
			Where("journal_entries.is_repaid = ?", true).
			Where("journal_entries.status IN ?", []string{"paid", "done"})

	case "inventory":
		return query.
			Where("journal_lines.debit > 0").
			Where("transaction_categories.debit_account_type = ?", "Asset").
			Where("transaction_categories.debit_category = ?", "Barang Dagangan").
			Where("transaction_categories.credit_account_type IN ?", []string{"Asset", "Liability"}).
			Where("transaction_categories.transaction_type = ?", "payout").
			Where("journal_entries.is_repaid = ?", true).
			Where("journal_entries.status IN ?", []string{"paid", "done"})

	case "receivable":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Revenue").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_entries.status = ?", "unpaid").
			Where("journal_entries.is_repaid = ?", false)

	case "receivable_history":
		return query.
			Where("journal_lines.debit > 0").
			Where("journal_lines.debit_account_type = ?", "Asset").
			Where("journal_lines.credit_account_type = ?", "Revenue").
			Where("journal_entries.transaction_type = ?", "payin").
			Where("journal_entries.status = ?", "done").
			Where("journal_entries.is_repaid = ?", true)

	default:
		return query
	}
}
