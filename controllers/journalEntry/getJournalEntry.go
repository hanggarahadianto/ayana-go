package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetJournalEntriesByCategory(c *gin.Context) {
	// Ambil query parameters dari request
	transactionCategoryID := c.DefaultQuery("transaction_category_id", "")
	companyID := c.DefaultQuery("company_id", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validasi parameter
	if transactionCategoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction category ID is required"})
		return
	}

	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	// Tentukan nilai default untuk page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var journalEntries []models.JournalEntry
	var total int64

	// Hitung total data berdasarkan transaction_category_id dan company_id
	if err := db.DB.Model(&models.JournalEntry{}).
		Where("transaction_category_id = ? AND company_id = ?", transactionCategoryID, companyID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count journal entries"})
		return
	}

	// Ambil data jurnal berdasarkan transaction_category_id dan company_id
	if err := db.DB.Preload("Lines").
		Where("transaction_category_id = ? AND company_id = ?", transactionCategoryID, companyID).
		Limit(limit).
		Offset(offset).
		Find(&journalEntries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch journal entries"})
		return
	}

	// Kirimkan response dengan data jurnal yang ditemukan
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   journalEntries,
	})
}
