package controller

import (
	"ayana/service"
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
	expenseStatus := c.DefaultQuery("status", "")
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	debitCategory := c.Query("debit_category")
	creditCategory := c.Query("credit_category")
	search := c.Query("search")

	if summaryOnlyStr != "true" && summaryOnlyStr != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter summary_only harus 'true' atau 'false'."})
		return
	}

	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	pagination := helper.GetPagination(c)

	params := service.ExpenseFilterParams{
		CompanyID:      companyID.String(),
		Pagination:     pagination,
		DateFilter:     dateFilter,
		ExpenseStatus:  expenseStatus,
		SummaryOnly:    summaryOnly,
		DebitCategory:  debitCategory,
		CreditCategory: creditCategory,
		Search:         search,
	}

	data, totalexpense, total, err := service.GetExpensesFromJournalLines(params)
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
