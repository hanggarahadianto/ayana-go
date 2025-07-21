package controller

import (
	lib "ayana/lib"
	expense "ayana/service/finance/expense"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetExpensesSummary(c *gin.Context) {

	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	accountType := "Expense"
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	debitCategory := c.Query("debit_category")
	creditCategory := c.Query("credit_category")
	search := c.Query("search")
	transactionType := c.DefaultQuery("transaction_type", "")
	expenseType := c.DefaultQuery("expense_type", "")

	if summaryOnlyStr != "true" && summaryOnlyStr != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter summary_only harus 'true' atau 'false'."})
		return
	}

	dateFilter, err := lib.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}
	sortBy := c.DefaultQuery("sort_by", "date_inputed") // default: date_inputed
	sortOrder := c.DefaultQuery("sort_order", "asc")    // default: asc

	pagination := lib.GetPagination(c)

	params := expense.ExpenseFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		AccountType:     accountType,
		TransactionType: transactionType,
		ExpenseType:     expenseType,
		SummaryOnly:     summaryOnly,
		DebitCategory:   debitCategory,
		CreditCategory:  creditCategory,
		Search:          search,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	// b, _ := json.MarshalIndent(params, "", "  ")
	// log.Println("Params:", string(b))

	data, totalexpense, total, err := expense.GetExpensesFromJournalLines(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"expenseList":   data,
			"total_expense": totalexpense,
			"page":          pagination.Page,
			"limit":         pagination.Limit,
			"total":         total,
		},
		"message": "Pengeluaran berhasil diambil",
		"status":  "sukses",
	})
}
