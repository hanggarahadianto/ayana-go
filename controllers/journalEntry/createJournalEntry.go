package controller

import (
	"net/http"

	"ayana/dto"
	"ayana/models"
	"ayana/service"

	"github.com/gin-gonic/gin"
)

func CreateJournalEntry(c *gin.Context) {
	var inputSingle models.JournalEntry
	var inputMultiple []models.JournalEntry

	// Coba parse multiple (array)
	if err := c.ShouldBindJSON(&inputMultiple); err == nil {
		results, err := service.ProcessMultipleJournalEntries(inputMultiple)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to process entries", "details": err.Error()})
			return
		}
		// Menggunakan DTO untuk hasil respons
		response := dto.MapToJournalEntryResponses(results)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
		return
	}

	// Kalau gagal, coba parse single
	if err := c.ShouldBindJSON(&inputSingle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input format", "details": err.Error()})
		return
	}

	result, err := service.ProcessSingleJournalEntry(inputSingle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to process entry", "details": err.Error()})
		return
	}

	// Menggunakan DTO untuk hasil respons
	response := dto.MapToJournalEntryResponse(result)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
