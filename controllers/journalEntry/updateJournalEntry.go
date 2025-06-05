package controller

import (
	"ayana/dto"
	"ayana/models"
	"ayana/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateJournalEntry(c *gin.Context) {

	var input models.JournalEntry
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	updatedJournal, err := service.UpdateSingleJournalEntry(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update journal entry",
			"details": err.Error(),
		})
		return
	}

	response := dto.MapToJournalEntryResponse(updatedJournal)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}
