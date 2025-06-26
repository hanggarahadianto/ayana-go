package helper

import (
	"fmt"
	"strings"
)

func BuildTypesenseExpenseTypeFilter(status string) string {
	var filters []string
	fmt.Println("🔥 Triggered BuildTypesenseExpenseTypeFilter with status = base")

	switch status {
	case "base":
		filters = append(filters,
			"transaction_type:=payout",
			"debit_account_type:=Expense",
			"credit_account_type:=Asset",
			"is_repaid:=true",
			"status:=[paid,done]",
		)
	}

	return strings.Join(filters, " && ")
}
