package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCashFlowById(c *gin.Context) {
	// Get the ID from URL params
	CashFlowId := c.Param("id")

	// Convert string ID to UUID
	cashFlowUUID, err := uuid.Parse(CashFlowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid CashFlow ID format",
		})
		return
	}

	// Define a variable to hold the CashFlow data
	var cashFlow models.CashFlow

	// Query the database with Preload to include related Goods
	result := db.DB.Preload("Good").Where("id = ?", cashFlowUUID).First(&cashFlow)

	// Check if the query resulted in an error
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "CashFlow not found",
		})
		return
	}

	// Return the cash flow data with related goods
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   cashFlow,
	})
}
