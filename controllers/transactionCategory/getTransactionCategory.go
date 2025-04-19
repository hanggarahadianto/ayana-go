package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTransactionCategory(c *gin.Context) {
	companyID := c.Query("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var transactions []models.TransactionCategory
	var total int64

	// Hitung total data
	if err := db.DB.Model(&models.TransactionCategory{}).
		Where("company_id = ?", companyID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transaction categories"})
		return
	}

	// Ambil data dengan relasi akun
	if err := db.DB.Preload("DebitAccount").
		Preload("CreditAccount").
		Where("company_id = ?", companyID).
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   transactions,
	})
}
