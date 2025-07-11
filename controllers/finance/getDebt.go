package controller

import (
	lib "ayana/lib"
	debt "ayana/service/finance/debt"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOutstandingDebts(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	accountType := "Debt" // Default account type for debts
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	if summaryOnlyStr != "true" && summaryOnlyStr != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter summary_only harus 'true' atau 'false'."})
		return
	}
	debitCategory := c.Query("debit_category")
	creditCategory := c.Query("credit_category")
	search := c.Query("search")

	debtType := c.DefaultQuery("debt_type", "")
	transactionType := c.DefaultQuery("transaction_type", "")

	dateFilter, err := lib.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}
	sortBy := c.DefaultQuery("sort_by", "date_inputed") // default: date_inputed
	sortOrder := c.DefaultQuery("sort_order", "asc")    // default: asc

	pagination := lib.GetPagination(c)

	params := debt.DebtFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		AccountType:     accountType,
		TransactionType: transactionType,
		DebtType:        debtType,
		SummaryOnly:     summaryOnly,
		DebitCategory:   debitCategory,
		CreditCategory:  creditCategory,
		Search:          search,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	data, totalDebt, total, err := debt.GetDebtsFromJournalLines(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"debtList":   data,
			"total_debt": totalDebt,
			"page":       pagination.Page,
			"limit":      pagination.Limit,
			"total":      total,
		},
		"message": "Hutang berhasil diambil",
		"status":  "sukses",
	})
}
