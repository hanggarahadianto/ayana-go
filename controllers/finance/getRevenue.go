package controller

import (
	"ayana/service"
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
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	category := c.Query("category")
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

	revenueType := c.DefaultQuery("revenue_type", "")
	transactionType := c.DefaultQuery("transaction_type", "")
	pagination := helper.GetPagination(c)
	params := service.RevenueFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		RevenueType:     revenueType,
		TransactionType: transactionType,
		SummaryOnly:     summaryOnly,
		Category:        category,
		Search:          search,
	}

	data, totalRevenue, total, err := service.GetRevenueFromJournalLines(params)
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
