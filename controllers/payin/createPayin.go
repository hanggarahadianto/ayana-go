package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePayIn(c *gin.Context) {
	var input models.JournalEntry

	// Validasi input dari frontend
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Siapkan ID & timestamps
	input.ID = uuid.New()
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	// Mulai transaksi
	tx := db.DB.Begin()

	// Simpan journal entry (header)
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan jurnal utama: " + err.Error()})
		return
	}

	// Persiapkan setiap line
	for i := range input.Lines {
		input.Lines[i].ID = uuid.New()
		input.Lines[i].JournalID = input.ID
		input.Lines[i].CreatedAt = time.Now()
		input.Lines[i].UpdatedAt = time.Now()
	}

	// Simpan journal lines (detail)
	if err := tx.Create(&input.Lines).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan jurnal detail: " + err.Error()})
		return
	}

	// Commit jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal commit transaksi: " + err.Error()})
		return
	}

	// Sukses
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": input})
}
