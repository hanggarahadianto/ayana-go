package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCashFlow(c *gin.Context) {
	var cashFlow models.CashFlow

	// Bind the incoming JSON to the cashFlow struct
	if err := c.ShouldBindJSON(&cashFlow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if the project ID exists in the Project table
	var project models.Project
	if err := db.DB.Where("id = ?", cashFlow.ProjectID).First(&project).Error; err != nil {
		// If no project is found, return an error message
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID does not exist"})
		} else {
			// Handle other potential errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking Project ID"})
		}
		return
	}

	// Start a transaction
	tx := db.DB.Begin()

	// Create the new CashFlow entry
	newCashFlow := models.CashFlow{
		ID:          uuid.New(),
		WeekNumber:  cashFlow.WeekNumber,
		CashIn:      cashFlow.CashIn,
		CashOut:     cashFlow.CashOut,
		Outstanding: cashFlow.Outstanding,
		ProjectID:   cashFlow.ProjectID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert the new CashFlow record into the database
	if err := tx.Create(&newCashFlow).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CashFlow"})
		return
	}

	// Check if Goods data is provided in the incoming request
	if len(cashFlow.Good) > 0 {
		// Prepare Goods entries with the new CashFlow ID
		var goods []models.Goods
		for _, good := range cashFlow.Good {
			good.ID = uuid.New()
			good.CashFlowId = newCashFlow.ID // Link the good to the new CashFlow
			goods = append(goods, good)
		}

		// Insert the Goods entries into the database
		if err := tx.Create(&goods).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Goods"})
			return
		}

		// Associate the created Goods with the new CashFlow
		newCashFlow.Good = goods
	}

	// Commit the transaction if all operations were successful
	tx.Commit()

	// Respond with success and return the created CashFlow data
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"cash_flow": newCashFlow,
		},
	})
}
