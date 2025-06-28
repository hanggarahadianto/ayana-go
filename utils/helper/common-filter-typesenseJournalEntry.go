package helper

import (
	"fmt"
	"strings"
	"time"
)

func BuildTypesenseFilter(
	companyID string,
	startDate, endDate *time.Time,
	accountType string,
	transactionType *string,
	Type *string,
	debitCategory string,
	creditCategory string,
) string {
	var filters []string

	// Pastikan companyID tidak kosong
	if companyID != "" {
		filters = append(filters, fmt.Sprintf("company_id:=%q", companyID))
	}

	// Asset type filter
	if Type != nil {
		fmt.Println("🔥 Triggered type =", *Type)

		switch accountType {
		case "Asset":
			if f := BuildTypesenseAssetTypeFilter(*Type); f != "" {
				filters = append(filters, f)
			}
		case "Expense":
			if f := BuildTypesenseExpenseTypeFilter(*Type); f != "" {
				filters = append(filters, f)
			}
		case "Debt":
			if f := BuildTypesenseDebtTypeFilter(*Type); f != "" {
				filters = append(filters, f)
			}
		}
	}
	if *transactionType != "" {
		filters = append(filters, fmt.Sprintf("transaction__type:=%q", *transactionType))
	}

	// Kategori debit dan kredit
	if debitCategory != "" {
		filters = append(filters, fmt.Sprintf("debit_category:=%q", debitCategory))
	}
	if creditCategory != "" {
		filters = append(filters, fmt.Sprintf("credit_category:=%q", creditCategory))
	}

	// Date range filter (timestamp dalam detik)
	if startDate != nil {
		filters = append(filters, fmt.Sprintf("date_inputed:>=%d", startDate.Unix()))
	}
	if endDate != nil {
		filters = append(filters, fmt.Sprintf("date_inputed:<=%d", endDate.Unix()))
	}

	// Gabungkan semua filter dengan " && "
	return strings.Join(filters, " && ")
}
