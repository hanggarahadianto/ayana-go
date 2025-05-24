package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DeleteJournalEntry(c *gin.Context) {
	// Ambil ID dari parameter URL
	idStr := c.Param("id")

	// Parse ID ke UUID
	journalEntryID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid journal entry ID"})
		return
	}

	// Cari journal entry berdasarkan ID
	var journalEntry models.JournalEntry
	err = db.DB.Preload("Lines").First(&journalEntry, "id = ?", journalEntryID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Journal entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find journal entry"})
		return
	}

	// Hapus relasi Lines terlebih dahulu (jika diperlukan)
	err = db.DB.Where("journal_id = ?", journalEntryID).Delete(&models.JournalLine{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete journal lines"})
		return
	}

	// Hapus journal entry
	err = db.DB.Delete(&journalEntry).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete journal entry"})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Journal entry deleted successfully",
	})
}
