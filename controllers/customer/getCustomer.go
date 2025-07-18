package controllers

import (
	"ayana/db"
	lib "ayana/lib"
	"ayana/models"
	customer "ayana/service/customer"

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
	status := c.Query("status") // ➕ Tambah ini
	summaryOnly := summaryOnlyStr == "true"
	sortBy := c.DefaultQuery("sort_by", "date_inputed")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	dateFilter, err := lib.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	selectStatus := c.Query("select_status")
	if selectStatus == "true" {
		var statuses []string
		if err := db.DB.Model(&models.Customer{}).
			Distinct("status").
			Pluck("status", &statuses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   statuses,
			"status": "success",
		})
		return
	}

	hasTestimonyStr := c.Query("has_testimony")
	var hasTestimony *bool
	if hasTestimonyStr != "" {
		v := hasTestimonyStr == "true"
		hasTestimony = &v
	}

	params := customer.CustomerFilterParams{
		CompanyID:    companyID,
		Pagination:   pagination,
		Search:       search,
		Status:       status, // ➕ Ini juga
		SummaryOnly:  summaryOnly,
		HasTestimony: hasTestimony,
		DateFilter:   dateFilter,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	data, total, err := customer.GetCustomersWithSearch(params)
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
