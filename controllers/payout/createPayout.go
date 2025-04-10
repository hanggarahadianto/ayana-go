package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreatePayout(c *gin.Context) {
	var payout models.Payout

	// Bind JSON input ke struct Payout
	if err := c.ShouldBindJSON(&payout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi apakah company_id ada di tabel Company
	var company models.Company
	if err := db.DB.First(&company, "id = ?", payout.CompanyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// Mulai transaksi database
	tx := db.DB.Begin()

	// Buat payout baru
	newPayout := models.Payout{

		Invoice:     payout.Invoice,
		Nominal:     payout.Nominal,
		DateInputed: payout.DateInputed,
		DueDate:     payout.DueDate,
		Note:        payout.Note,
		Status:      payout.Status,
		CompanyID:   payout.CompanyID,
		Category:    payout.Category,
		Mitra:       payout.Mitra,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Simpan payout ke database
	if err := tx.Create(&newPayout).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Payout"})
		return
	}

	// Commit transaksi
	tx.Commit()

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newPayout,
	})
}
