// helper/typesense_filter.go
package helper

import (
	"strings"
)

func BuildTypesenseAssetTypeFilter(assetType string) string {
	var filters []string

	switch assetType {
	case "cashin":
		filters = append(filters,
			"transaction_type:=payin",
			"debit_account_type:=Asset",
		)

	case "cashout":
		filters = append(filters,
			"transaction_type:=payout",
			"credit_account_type:=Asset",
			"status:=[paid,done]",
			"is_repaid:=true",
		)

	case "fixed_asset":
		filters = append(filters,
			"transaction_type:=payout",
			"debit_account_type:=Asset",
			"debit_category:=Aset Tetap",
			"credit_account_type:=[Asset,Liability]",
			"status:=[paid,done]",
			"is_repaid:=true",
		)

	case "inventory":
		filters = append(filters,
			"transaction_type:=payout",
			"debit_account_type:=Asset",
			"debit_category:=Barang Dagangan",
			"credit_account_type:=[Asset,Liability]",
			"status:=[paid,done]",
			"is_repaid:=true",
		)

	case "receivable":
		filters = append(filters,
			"transaction_type:=payin",
			"debit_account_type:=Asset",
			"credit_account_type:=Revenue",
			"status:=unpaid",
			"is_repaid:=false",
		)

	case "receivable_history":
		filters = append(filters,
			"transaction_type:=payin",
			"debit_account_type:=Asset",
			"credit_account_type:=Revenue",
			"status:=done",
			"is_repaid:=true",
		)
	}

	return strings.Join(filters, " && ")
}
