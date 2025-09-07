package controller

import (
	"ayana/dto"
	"ayana/models"
	journalEntry "ayana/service/journalEntry"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateReversedJournalEntry(c *gin.Context) {

	var raw json.RawMessage

	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Cek apakah input berupa array atau object
	if raw[0] == '[' {
		var inputMultiple []models.JournalEntry
		if err := json.Unmarshal(raw, &inputMultiple); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid array format", "details": err.Error()})
			return
		}

		results, err := journalEntry.ProcessReverseJournalEntry(inputMultiple)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to process entries", "details": err.Error()})
			return
		}

		response := dto.MapToJournalEntryResponses(results)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
		return
	}

}
