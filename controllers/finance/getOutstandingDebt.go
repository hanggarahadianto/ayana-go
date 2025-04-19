package controller

import (
	"ayana/db"
	"ayana/models"
	validationEnum "ayana/utils/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOutstandingDebts(c *gin.Context) {

	companyID := c.Query("company_id")
	transactionType := c.Query("transaction_type")
	status := c.Query("status")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}
	if transactionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TransactionType is required"})
		return
	}
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	if !validationEnum.IsValidStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if !validationEnum.IsValidTransactionType(transactionType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction_type"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	// var entries []JournalEntry
	var entries []models.JournalEntry
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := db.DB.
		Where("company_id = ? AND transaction_type = ? AND status = ? AND is_repaid = false", companyID, transactionType, status)

	// Count total matching records
	if err := query.Model(&models.JournalEntry{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to count data"})
		return
	}

	if err := db.DB.
		Preload("Lines").
		Preload("Lines.Account").
		Preload("TransactionCategory").
		Preload("TransactionCategory.DebitAccount").
		Preload("TransactionCategory.CreditAccount").
		Where("company_id = ? AND transaction_type = ? AND status = ? AND is_repaid = false", companyID, transactionType, status).
		Limit(limit).
		Offset(offset).
		Order("due_date ASC").
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   entries,
	})
}
