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
	// Bind JSON input ke struct Payout
	var updatedPayout models.Payout
	if err := c.ShouldBindJSON(&updatedPayout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi apakah ID valid
	if updatedPayout.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payout ID"})
		return
	}

	// Cek apakah Payout ada di database berdasarkan ID dan CompanyID
	var payout models.Payout
	if err := db.DB.Where("id = ? AND company_id = ?", updatedPayout.ID, updatedPayout.CompanyID).First(&payout).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found"})
		return
	}

	// Mulai transaksi database
	tx := db.DB.Begin()

	// Update data payout
	payout.Invoice = updatedPayout.Invoice
	payout.Nominal = updatedPayout.Nominal
	payout.DateInputed = updatedPayout.DateInputed
	payout.Note = updatedPayout.Note
	payout.UpdatedAt = time.Now()

	if err := tx.Save(&payout).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Payout"})
		return
	}

	// Commit transaksi
	tx.Commit()

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   payout,
	})
}
