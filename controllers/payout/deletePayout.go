package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeletePayout(c *gin.Context) {
	// Ambil ID dari parameter URL
	payoutID := c.Param("id")

	// Parse ID ke UUID
	parsedID, err := uuid.Parse(payoutID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payout ID"})
		return
	}

	// Cek apakah Payout ada di database
	var payout models.Payout
	if err := db.DB.First(&payout, "id = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found"})
		return
	}

	// Mulai transaksi database
	tx := db.DB.Begin()

	// Hapus payout dari database
	if err := tx.Delete(&payout).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Payout"})
		return
	}

	// Commit transaksi
	tx.Commit()

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Payout deleted successfully",
	})
}
