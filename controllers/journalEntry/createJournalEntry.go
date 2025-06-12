package controller

import (
	"encoding/json"
	"net/http"

	"ayana/dto"
	"ayana/models"
	"ayana/service"

	"github.com/gin-gonic/gin"
)

func CreateJournalEntry(c *gin.Context) {
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

		results, err := service.ProcessMultipleJournalEntries(inputMultiple)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to process entries", "details": err.Error()})
			return
		}

		response := dto.MapToJournalEntryResponses(results)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
		return
	}

	// Jika bukan array, asumsikan single object
	var inputSingle models.JournalEntry
	if err := json.Unmarshal(raw, &inputSingle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid object format", "details": err.Error()})
		return
	}

	result, err := service.ProcessSingleJournalEntry(inputSingle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to process entry", "details": err.Error()})
		return
	}

	response := dto.MapToJournalEntryResponse(result)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
