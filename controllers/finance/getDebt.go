package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOutstandingDebts(c *gin.Context) {

	companyIDStr := c.Query("company_id")

	// Mendapatkan parameter summary_only
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")

	// Menangani logika untuk summary_only
	summaryOnly := false
	if summaryOnlyStr == "true" {
		summaryOnly = true
	} else if summaryOnlyStr != "false" {
		// Menangani kasus di mana nilai selain "true" atau "false" diberikan
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter summary_only harus bernilai 'true' atau 'false'."})
		return
	}

	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	// Ambil date filter
	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	pagination := helper.GetPagination(c)

	params := service.DebtFilterParams{
		CompanyID: companyID.String(),

		Pagination:  pagination,
		DateFilter:  dateFilter,
		SummaryOnly: summaryOnly, // Menambahkan summaryOnly pada parameter
	}

	data, totalDebt, total, err := service.GetOutstandingDebts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"debtList":              data,
			"total_outstandingDebt": totalDebt,
			"page":                  pagination.Page,
			"limit":                 pagination.Limit,
			"total":                 total,
		},
		"message": "Outstanding debts retrieved successfully",
		"status":  "success",
	})
}
