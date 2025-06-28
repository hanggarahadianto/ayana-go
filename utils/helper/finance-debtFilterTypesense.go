package helper

import (
	"fmt"
	"strings"
)

func BuildTypesenseDebtTypeFilter(debtType string) string {
	var filters []string

	switch debtType {
	case "going":
		fmt.Println("🔥 Triggered BuildTypesenseDebtTypeFilter with status = going")
		filters = append(filters,
			"credit:>0",
			"is_repaid:=false",
			"status:=unpaid",
			"credit_account_type:=Liability",
			"debit_account_type:!=Revenue",
		)

	case "done":
		fmt.Println("🔥 Triggered BuildTypesenseDebtTypeFilter with status = done")
		filters = append(filters,
			"debit:>0",
			"is_repaid:=true",
			"status:=done",
			"debit_account_type:=Liability",
			"credit_account_type:=Asset",
		)
	}

	return strings.Join(filters, " && ")
}
