package controllers

import (
	"ayana/service"
	"ayana/utils/helper"

	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	pagination := helper.GetPagination(c)
	if !helper.ValidatePagination(pagination, c) {
		return
	}

	search := c.Query("search")
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"

	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	params := service.CustomerFilterParams{
		Pagination:  pagination,
		Search:      search,
		SummaryOnly: summaryOnly,
		DateFilter:  dateFilter,
	}

	data, total, err := service.GetCustomersWithSearch(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"customerList":   data,
			"total_customer": total,
			"page":           pagination.Page,
			"limit":          pagination.Limit,
			"total":          total,
		},
		"message": "Data customer berhasil diambil",
		"status":  "sukses",
	})
}
