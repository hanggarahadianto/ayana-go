package helper

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidateAndParseTransactionCategoryID(transactionCategoryIDStr string, c *gin.Context) (uuid.UUID, bool) {
	if transactionCategoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction Category ID is required"})
		return uuid.Nil, false
	}

	transactionCategoryID, err := uuid.Parse(transactionCategoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for Transaction Category ID"})
		return uuid.Nil, false
	}

	var category models.TransactionCategory
	if err := db.DB.First(&category, "id = ?", transactionCategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Transaction Category ID"})
		return uuid.Nil, false
	}

	return transactionCategoryID, true
}
