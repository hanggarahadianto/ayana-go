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

func CreateWeeklyProgress(c *gin.Context) {
	var weeklyProgress models.WeeklyProgress

	// Bind the incoming JSON to the weeklyProgress struct
	if err := c.ShouldBindJSON(&weeklyProgress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if the project ID exists in the Project table
	var project models.Project
	if err := db.DB.Where("id = ?", weeklyProgress.ProjectID).First(&project).Error; err != nil {
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

	// Create the WeeklyProgress entry
	newWeeklyProgress := models.WeeklyProgress{
		ID:             uuid.New(),
		WeekNumber:     weeklyProgress.WeekNumber,
		Percentage:     weeklyProgress.Percentage,
		AmountMaterial: weeklyProgress.AmountMaterial,
		AmountWorker:   weeklyProgress.AmountWorker,
		ProjectID:      weeklyProgress.ProjectID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Insert the WeeklyProgress
	if err := tx.Create(&newWeeklyProgress).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create WeeklyProgress"})
		return
	}

	// Prepare Material entries with the WeeklyProgress ID
	var materials []models.Material
	for _, material := range weeklyProgress.Material {
		material.ID = uuid.New()
		material.WeeklyProgressIdMaterial = newWeeklyProgress.ID
		materials = append(materials, material)
	}

	// Insert Material entries
	if err := tx.Create(&materials).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Materials"})
		return
	}

	// Prepare Worker entries with the WeeklyProgress ID
	var workers []models.Worker
	for _, worker := range weeklyProgress.Worker {
		worker.ID = uuid.New()
		worker.WeeklyProgressIdWorker = newWeeklyProgress.ID
		workers = append(workers, worker)
	}

	// Insert Worker entries
	if err := tx.Create(&workers).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Workers"})
		return
	}

	// Commit the transaction if all operations were successful
	tx.Commit()

	// Add materials and workers data into the weekly_progress object
	newWeeklyProgress.Material = materials
	newWeeklyProgress.Worker = workers

	// Respond with success and return the created data
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"weekly_progress": newWeeklyProgress,
		},
	})
}
