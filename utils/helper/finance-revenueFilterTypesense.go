package helper

import (
	"fmt"
	"strings"
)

func BuildTypesenseRevenueTypeFilter(revenueType string) string {
	var filters []string

	if revenueType == "realized" {
		fmt.Println("🔥 Triggered BuildTypesenseExpenseTypeFilter with status = realized")
		filters = append(filters,
			"transaction_type:=payin",
			"debit_account_type:=Asset",
			"credit_account_type:=Revenue",
			"is_repaid:=true",
			"status:=[paid,done]",
		)
	}

	return strings.Join(filters, " && ")
}
