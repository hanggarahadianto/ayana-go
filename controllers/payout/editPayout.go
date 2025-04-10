package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func EditPayout(c *gin.Context) {
	var input models.Payout

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payout ID"})
		return
	}

	// Ambil payout berdasarkan ID dan company_id
	var existing models.Payout
	if err := db.DB.Where("id = ? AND company_id = ?", input.ID, input.CompanyID).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found"})
		return
	}

	// Update field jika nilainya tidak kosong atau nol
	if input.Invoice != "" {
		existing.Invoice = input.Invoice
	}
	if input.Nominal != 0 {
		existing.Nominal = input.Nominal
	}
	if input.DateInputed != nil {
		existing.DateInputed = input.DateInputed
	}
	if input.DueDate != nil {
		existing.DueDate = input.DueDate
	}
	if input.PaymentDate != nil {
		existing.PaymentDate = input.PaymentDate
	}

	if input.Note != "" {
		existing.Note = input.Note
	}
	if input.Status != "" {
		existing.Status = input.Status
	}

	existing.UpdatedAt = time.Now()

	// Simpan perubahan
	if err := db.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existing,
	})
}
