package controller

import (
	"ayana/dto"
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransactionCategory(c *gin.Context) {

	pagination := helper.GetPagination(c)

	filterParams := service.TransactionCategoryFilterParams{
		CompanyID:       c.Query("company_id"),
		TransactionType: c.Query("transaction_type"),
		Category:        c.Query("category"),
		Status:          c.Query("status"),
		All:             c.Query("all") == "true",
		Pagination:      pagination,
	}

	var data []dto.TransactionCategoryResponse
	var total int64
	var err error

	if pagination.Limit == 0 && !filterParams.All {
		data, err := service.GetTransactionCategoriesWithoutPagination(filterParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
		return
	}

	data, total, err = service.GetTransactionCategories(filterParams)
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
