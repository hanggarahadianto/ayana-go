package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateTransactionCategory(c *gin.Context) {
	// Validasi ID param
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Bind JSON ke struct
	var input models.TransactionCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Cari data existing berdasarkan ID
	var existing models.TransactionCategory
	if err := db.DB.First(&existing, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction category not found"})
		return
	}

	// Update field yang diizinkan
	existing.Name = input.Name
	existing.TransactionType = input.TransactionType

	existing.DebitAccountID = input.DebitAccountID
	existing.CreditAccountID = input.CreditAccountID

	existing.DebitAccountType = input.DebitAccountType
	existing.CreditAccountType = input.CreditAccountType

	existing.DebitCategory = input.DebitCategory
	existing.CreditCategory = input.CreditCategory

	existing.TransactionLabel = input.TransactionLabel

	existing.Description = input.Description
	existing.Status = input.Status
	existing.CompanyID = input.CompanyID

	// Simpan perubahan
	if err := db.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction category"})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existing,
	})
}
