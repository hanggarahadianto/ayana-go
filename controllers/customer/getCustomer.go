package controllers

import (
	lib "ayana/lib"
	"ayana/service"

	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	pagination := lib.GetPagination(c)
	if !lib.ValidatePagination(pagination, c) {
		return
	}

	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id wajib diisi"})
		return
	}

	search := c.Query("search")
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	sortBy := c.DefaultQuery("sort_by", "date_inputed")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	dateFilter, err := lib.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	params := service.CustomerFilterParams{
		CompanyID:   companyID,
		Pagination:  pagination,
		Search:      search,
		SummaryOnly: summaryOnly,
		DateFilter:  dateFilter,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}

	data, total, err := service.GetCustomersWithSearch(params)
	if err != nil {
		log.Println("🔴 GetCustomers error:", err) // Tambahkan log ini untuk debugging
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
