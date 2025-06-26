package controller

import (
	lib "ayana/lib"
	revenue "ayana/service/finance/revenue"
	"ayana/utils/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRevenueSummary(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	accountType := "revenue" // Default account type for revenue
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

	revenueStatus := c.DefaultQuery("revenue_type", "")
	transactionType := c.DefaultQuery("transaction_type", "")
	pagination := lib.GetPagination(c)
	params := revenue.RevenueFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		AccountType:     accountType,
		Status:          revenueStatus,
		TransactionType: transactionType,
		SummaryOnly:     summaryOnly,
		DebitCategory:   debitCategory,
		CreditCategory:  creditCategory,
		Search:          search,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	data, totalRevenue, total, err := revenue.GetRevenueFromJournalLines(params)
	if err != nil {
		log.Printf("GetRevenueFromJournalLines error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aset"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"revenueList":   data,
			"total_revenue": totalRevenue,
			"page":          pagination.Page,
			"limit":         pagination.Limit,
			"total":         total,
		},
		"message": "Pendapatan berhasil diambil",
		"status":  "sukses",
	})
}
