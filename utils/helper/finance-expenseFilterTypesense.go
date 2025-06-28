package helper

import (
	"fmt"
	"strings"
)

func BuildTypesenseExpenseTypeFilter(expenseType string) string {
	var filters []string

	if expenseType == "base" {
		fmt.Println("🔥 Triggered BuildTypesenseExpenseTypeFilter with status = base")
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
