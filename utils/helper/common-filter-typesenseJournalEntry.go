package helper

import (
	"fmt"
	"strings"
	"time"
)

func BuildTypesenseFilter(
	companyID string,
	debitCategory string,
	creditCategory string,
	startDate, endDate *time.Time,
	status *string,
	assetType *string,
) string {
	var filters []string

	// Pastikan companyID tidak kosong
	if companyID != "" {
		filters = append(filters, fmt.Sprintf("company_id:=%q", companyID))
	}

	// Asset type filter
	if assetType != nil {
		if f := BuildTypesenseAssetTypeFilter(*assetType); f != "" {
			filters = append(filters, f)
		}
	}

	// Status filter
	if status != nil && *status != "" {
		filters = append(filters, fmt.Sprintf("status:=%q", *status))
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
