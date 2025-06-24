package controller

import (
	"ayana/db"
	"ayana/models"
	service "ayana/service/journalEntry"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DeleteJournalEntries(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"` // terima array string UUID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validasi dan konversi ke UUID
	var uuidList []uuid.UUID
	for _, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID in list", "id": idStr})
			return
		}
		uuidList = append(uuidList, id)
	}

	// Mulai transaksi
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for _, journalEntryID := range uuidList {
			var journalEntry models.JournalEntry

			if err := tx.Preload("Lines").First(&journalEntry, "id = ?", journalEntryID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					continue // skip ID yang tidak ditemukan
				}
				return err
			}

			if err := tx.Where("journal_id = ?", journalEntryID).Delete(&models.JournalLine{}).Error; err != nil {
				return err
			}

			if err := tx.Delete(&journalEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete journal entries", "details": err.Error()})
		return
	}

	// Hapus juga dari Typesense (optional)
	for _, id := range uuidList {
		if err := service.DeleteJournalEntryFromTypesense(c.Request.Context(), id.String()); err != nil {
			// Optional: tangani error ini jika penting
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Journal entries deleted successfully",
	})
}
