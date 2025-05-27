package controller

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateReversedJournalEntry(c *gin.Context) {
	var inputEntries []models.JournalEntry

	// Parsing input JSON array
	if err := c.ShouldBindJSON(&inputEntries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}
	if len(inputEntries) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "No journal entry provided",
		})
		return
	}

	firstEntry := inputEntries[0]

	// ✅ Validasi first entry
	if err := helper.ValidateJournalEntry(firstEntry.ID, firstEntry.Transaction_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// ✅ Update first entry saja
	if err := db.DB.Model(&models.JournalEntry{}).
		Where("id = ?", firstEntry.ID).
		Updates(map[string]interface{}{
			"invoice":   firstEntry.Invoice,
			"status":    "paid",
			"is_repaid": true,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update journal entry",
			"details": err.Error(),
		})
		return
	}

	// ✅ Array 2: Create ONLY entries starting from index 1
	var newEntries []models.JournalEntry
	if len(inputEntries) > 1 {
		for _, entry := range inputEntries[1:] { // MULAI dari index ke-1
			newEntry := models.JournalEntry{
				Amount:                entry.Amount,
				Transaction_ID:        entry.Transaction_ID,
				CompanyID:             entry.CompanyID,
				DateInputed:           entry.DateInputed,
				Description:           entry.Description,
				DueDate:               entry.DueDate,
				Installment:           entry.Installment,
				Invoice:               entry.Invoice,
				IsRepaid:              entry.IsRepaid,
				Note:                  entry.Note,
				Partner:               entry.Partner,
				Status:                entry.Status,
				TransactionCategoryID: entry.TransactionCategoryID,
				TransactionType:       entry.TransactionType,
			}
			newEntries = append(newEntries, newEntry)
		}
	}

	// ✅ Hanya create kalau ada newEntries
	if len(newEntries) > 0 {
		result, err := service.ProcessMultipleJournalEntries(newEntries)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to process entries",
				"details": err.Error(),
			})
			return
		}

		// Response pakai DTO
		response := dto.MapToJournalEntryResponseList(result)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
		return
	}

	// ✅ Kalau tidak ada entry baru, tetap success
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "First journal entry updated successfully"})
}
