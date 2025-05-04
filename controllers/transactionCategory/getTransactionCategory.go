package controller

import (
	"ayana/dto"
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransactionCategory(c *gin.Context) {
	// Mendapatkan parameter untuk pagination
	pagination := helper.GetPagination(c)
	companyID := c.Query("company_id")
	transactionType := c.Query("transaction_type")
	category := c.Query("category")
	status := c.Query("status")

	// Menyusun filter parameter untuk query
	params := service.TransactionCategoryFilterParams{
		CompanyID:       companyID,
		TransactionType: transactionType,
		Category:        category,
		Status:          status,
		Pagination:      pagination,
	}

	// Panggil service untuk mendapatkan data berdasarkan parameter
	var data []dto.TransactionCategoryResponse
	var total int64
	var err error

	// Jika pagination limit 0, berarti kita ingin ambil semua data
	if pagination.Limit == 0 {
		data, err := service.GetTransactionCategoriesWithoutPagination(companyID, transactionType, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   data,
		})

	} else {
		data, total, err = service.GetTransactionCategories(params)
	}

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
