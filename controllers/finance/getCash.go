package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCashSummary(c *gin.Context) {
	// Step 1. Validasi company_id
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	// Step 2. Ambil date filter
	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 3. Kalau summary_only = true -> hitung total cash masuk saja
	if c.Query("summary_only") == "true" {
		totalCash, err := service.GetCashSummaryOnly(companyID.String(), dateFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"total_cash": totalCash,
			},
			"message": "Cash summary retrieved successfully",
			"status":  "success",
		})
		return
	}

	// Step 4. Get Pagination
	pagination := helper.GetPagination(c)

	// Step 5. Ambil data cash list
	cashList, total, count, err := service.GetCash(service.CashFilterParams{
		CompanyID:  companyID.String(),
		DateFilter: dateFilter,
		Pagination: pagination,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 6. Response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"cashList": cashList,
			"total":    total,
			"count":    count,
			"page":     pagination.Page,
			"limit":    pagination.Limit,
		},
		"message": "Cash summary retrieved successfully",
		"status":  "success",
	})
}
