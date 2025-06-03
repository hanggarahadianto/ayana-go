package controller

import (
	"ayana/db"
	"ayana/models"
	"ayana/service"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DeleteJournalEntry(c *gin.Context) {
	idStr := c.Param("id")

	journalEntryID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid journal entry ID"})
		return
	}

	// Mulai transaction
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Ambil journal entry
		var journalEntry models.JournalEntry
		if err := tx.Preload("Lines").First(&journalEntry, "id = ?", journalEntryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Journal entry not found"})
				return err
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find journal entry"})
			return err
		}

		// Hapus relasi Lines
		if err := tx.Where("journal_id = ?", journalEntryID).Delete(&models.JournalLine{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete journal lines"})
			return err
		}

		// Hapus journal entry
		if err := tx.Delete(&journalEntry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete journal entry"})
			return err
		}

		// Hapus dari Typesense (tidak ikut rollback, logika luar DB)
		if err := service.DeleteJournalEntryFromTypesense(c.Request.Context(), journalEntryID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from Typesense"})
			return err
		}

		return nil
	})

	// Jika transaksi gagal (sudah diproses di dalam)
	if err != nil {
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Journal entry deleted successfully",
	})
}
