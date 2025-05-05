// utils/helper.go

package helper

import (
	"strings"

	"gorm.io/gorm"
)

func ApplyTransactionFilters(tx *gorm.DB, transactionType, category, status string) *gorm.DB {
	if transactionType != "" {
		tx = tx.Where("LOWER(transaction_categories.transaction_type) = ?", strings.ToLower(transactionType))
	}

	if category != "" {
		tx = tx.Where("transaction_categories.category ILIKE ?", "%"+category+"%")
	}

	// ➕ Status-based logic
	if status == "unpaid" {
		// Uang keluar tapi belum dibayar → hutang
		tx = tx.
			// Where("transaction_categories.debit_account_type = ?", "Asset").
			Where("transaction_categories.credit_account_type = ?", "Liability")

	} else if status == "paid" {
		tx = tx.
			// Where("transaction_categories.debit_account_type = ?", "Asset").
			Where("transaction_categories.credit_account_type = ?", "Asset")
	}

	return tx
}
