package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOutstandingDebts(c *gin.Context) {

	status := c.Query("status")

	companyIDStr := c.Query("company_id")

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

	// Jika hanya ingin summary
	if c.Query("summary_only") == "true" {
		totalOutstandingDebt, err := service.GetOutstandingDebtSummaryOnly(companyID.String(), dateFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate summary total"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"total_outstandingDebt": totalOutstandingDebt,
			},
			"message": "Outstanding debt summary retrieved successfully",
			"status":  "success",
		})
		return
	}

	// Jika ingin data list + summary
	pagination := helper.GetPagination(c)

	params := service.DebtFilterParams{
		CompanyID:  companyID.String(),
		Status:     status,
		Pagination: pagination,
		DateFilter: dateFilter,
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
