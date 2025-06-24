package controller

import (
	lib "ayana/lib"
	"ayana/service"
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

	equityType := c.DefaultQuery("equity_type", "")
	transactionType := c.DefaultQuery("transaction_type", "")
	pagination := lib.GetPagination(c)
	params := service.EquityFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		EquityType:      equityType,
		TransactionType: transactionType,
		SummaryOnly:     summaryOnly,
		DebitCategory:   debitCategory,
		CreditCategory:  creditCategory,
		Search:          search,
	}

	data, totalEquity, total, err := service.GetEquityFromJournalLines(params)
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
