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
	transactionType := c.Query("transaction_type") // payin / payout

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

	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
		Where("transaction_categories.company_id = ?", companyID)

	// ✅ Filter berdasarkan jenis transaksi
	if transactionType == "payin" {
		// Payin → Filter debit account bertipe Asset
		tx = tx.Where("LOWER(debit.type) LIKE ?", "asset%")
	} else if transactionType == "payout" {
		// Payout → Join credit account dan filter tipe Asset
		tx = tx.Joins("JOIN accounts AS credit ON credit.id = transaction_categories.credit_account_id").
			Where("LOWER(credit.type) LIKE ?", "asset%")
	}

	// Hitung total
	if err := tx.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transaction categories"})
		return
	}

	// Ambil data dengan preload account
	if err := tx.Preload("DebitAccount").
		Preload("CreditAccount").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories"})
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   transactions,
	})
}
