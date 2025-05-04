package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTransactionCategory handler untuk mendapatkan kategori transaksi dengan filter
func GetTransactionCategory(c *gin.Context) {
	// Mendapatkan parameter untuk pagination
	pagination := helper.GetPagination(c)
	companyID := c.Query("company_id")
	transactionType := c.Query("transaction_type")

	category := c.Query("category") // ✅ Tambahan

	// Menyusun filter parameter untuk query
	params := service.TransactionCategoryFilterParams{
		CompanyID:       companyID,
		TransactionType: transactionType,
		Category:        category,

		Pagination: pagination,
	}

	// Panggil service untuk mendapatkan data berdasarkan parameter
	data, total, err := service.GetTransactionCategories(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories"})
		return
	}

	// Kirimkan response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   pagination.Page,
		"limit":  pagination.Limit,
		"total":  total,
		"data":   data,
	})
}
