package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateTransactionCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var input models.TransactionCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var existing models.TransactionCategory
	if err := db.DB.First(&existing, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction category not found"})
		return
	}

	// Update only allowed fields
	existing.CompanyID = input.CompanyID
	existing.TransactionType = input.TransactionType
	existing.Category = input.Category
	// existing.Status = input.Status

	if err := db.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existing,
	})
}
