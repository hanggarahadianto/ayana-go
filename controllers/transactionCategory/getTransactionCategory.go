package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransactionCategory(c *gin.Context) {
	pagination := helper.GetPagination(c)

	all := c.Query("all") == "true"
	selectOnly := c.Query("select") == "true"
	selectByCategory := c.Query("select_by_category") == "true" // ✅ Tambahan

	filterParams := service.TransactionCategoryFilterParams{
		CompanyID:         c.Query("company_id"),
		TransactionType:   c.Query("transaction_type"),
		DebitCategory:     c.Query("debit_category"),
		CreditCategory:    c.Query("credit_category"),
		Status:            c.Query("status"),
		DebitAccountType:  c.Query("debit_account_type"),  // ➕ Tambahan
		CreditAccountType: c.Query("credit_account_type"), // ➕ Tambahan
		All:               all,
		Pagination:        pagination,
		SelectOnly:        selectOnly, // tambah field baru di struct params

	}

	var data interface{}
	var total int64
	var err error

	if all {
		data, err = service.GetTransactionCategoriesAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all transaction categories"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
		return
	}
	if selectByCategory {
		// ✅ Tangani select kategori unik
		if filterParams.CompanyID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "company_id filter is required for selectByCategory=true"})
			return
		}

		data, message, err := service.GetUniqueCategories(filterParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unique categories"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": message,
			"data":   data,
		})
		return
	}

	if selectOnly {
		// Select = true, wajib ada filter, tanpa paginasi
		if filterParams.CompanyID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "company_id filter is required for select=true"})
			return
		}
		data, err = service.GetTransactionCategoriesForSelect(filterParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories for select"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
		return
	}

	// Default: paginasi dengan filter
	data, total, err = service.GetTransactionCategoriesWithPagination(filterParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   pagination.Page,
		"limit":  pagination.Limit,
		"total":  total,
		"data":   data,
	})
}
