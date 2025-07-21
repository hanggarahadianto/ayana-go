package helper

import (
	"fmt"
	"strings"
)

func BuildTypesenseEquityTypeFilter(equityType string) string {
	var filters []string

	switch equityType {
	case "setor":
		fmt.Println("🔥 Triggered BuildTypesenseEquityTypeFilter with type = setor")
		filters = append(filters,
			"transaction_type:=payin",
			"debit_account_type:=Asset",
			"credit_account_type:=Equity",
			"is_repaid:=true",
			"status:=[paid,done]",
		)
	case "tarik":
		fmt.Println("🔥 Triggered BuildTypesenseEquityTypeFilter with type = tarik")
		filters = append(filters,
			"transaction_type:=payout",
			"debit_account_type:=Equity",
			"credit_account_type:=Asset",
		)
	default:
		fmt.Printf("⚠️ Unknown equityType: %s\n", equityType)
	}

	filterString := strings.Join(filters, " && ")
	fmt.Println("🧪 Final Equity Filter:", filterString) // <--- LOG FINAL FILTER
	return filterString
}
