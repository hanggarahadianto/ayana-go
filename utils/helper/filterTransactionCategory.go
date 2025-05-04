// utils/helper.go

package helper

import (
	"strings"

	"gorm.io/gorm"
)

func ApplyTransactionFilters(tx *gorm.DB, transactionType, category string) *gorm.DB {
	if transactionType != "" {
		tx = tx.Where("LOWER(transaction_categories.transaction_type) = ?", strings.ToLower(transactionType))
	}

	if category != "" {
		// ✅ pakai ILIKE agar fleksibel (dan case-insensitive)
		tx = tx.Where("transaction_categories.category ILIKE ?", "%"+category+"%")
	}

	return tx
}
