package controller

import (
	lib "ayana/lib"
	equity "ayana/service/finance/equity"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEquitySummary(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	accountType := "Equity"
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	debitCategory := c.Query("debit_category")
	creditCategory := c.Query("credit_category")
	search := c.Query("search")

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

	equityType := c.DefaultQuery("equity_type", "")
	transactionType := c.DefaultQuery("transaction_type", "")
	pagination := lib.GetPagination(c)
	params := equity.EquityFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		AccountType:     accountType,
		EquityType:      equityType,
		TransactionType: transactionType,
		SummaryOnly:     summaryOnly,
		DebitCategory:   debitCategory,
		CreditCategory:  creditCategory,
		Search:          search,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	data, totalEquity, total, err := equity.GetEquityFromJournalLines(params)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data modal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"equityList":   data,
			"total_equity": totalEquity,
			"page":         pagination.Page,
			"limit":        pagination.Limit,
			"total":        total,
		},
		"message": "Modal berhasil diambil",
		"status":  "sukses",
	})
}
