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

func EditWeeklyProgress(c *gin.Context) {
	var weeklyProgress models.WeeklyProgress

	// Bind the incoming JSON to the weeklyProgress struct
	if err := c.ShouldBindJSON(&weeklyProgress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if the project ID exists in the Project table
	var project models.Project
	if err := db.DB.Where("id = ?", weeklyProgress.ProjectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking Project ID"})
		}
		return
	}

	// Fetch existing WeeklyProgress record, preload Material and Worker
	var existingWeeklyProgress models.WeeklyProgress
	if err := db.DB.Where("id = ?", weeklyProgress.ID).Preload("Material").Preload("Worker").First(&existingWeeklyProgress).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Weekly Progress ID does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Weekly Progress"})
		}
		return
	}

	// Start a transaction
	tx := db.DB.Begin()

	// Update WeeklyProgress fields
	existingWeeklyProgress.WeekNumber = weeklyProgress.WeekNumber
	existingWeeklyProgress.Percentage = weeklyProgress.Percentage
	existingWeeklyProgress.AmountMaterial = weeklyProgress.AmountMaterial
	existingWeeklyProgress.AmountWorker = weeklyProgress.AmountWorker
	existingWeeklyProgress.ProjectID = weeklyProgress.ProjectID
	existingWeeklyProgress.UpdatedAt = time.Now()

	// Save updated WeeklyProgress record
	if err := tx.Save(&existingWeeklyProgress).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update WeeklyProgress"})
		return
	}

	var workerIDsToDelete []uuid.UUID
	for _, existingWorker := range existingWeeklyProgress.Worker {
		found := false
		for _, updatedWorker := range weeklyProgress.Worker {
			if existingWorker.ID == updatedWorker.ID {
				found = true
				break
			}
		}
		if !found {
			// Mark for deletion
			workerIDsToDelete = append(workerIDsToDelete, existingWorker.ID)
		}
	}

	if len(workerIDsToDelete) > 0 {
		if err := tx.Where("id IN ?", workerIDsToDelete).Delete(&models.Worker{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Workers"})
			return
		}
	}

	// Handle Workers: Update existing or add new workers
	for _, worker := range weeklyProgress.Worker {
		if worker.ID != uuid.Nil {
			// Update existing worker
			if err := tx.Model(&models.Worker{}).Where("id = ?", worker.ID).Updates(worker).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Worker"})
				return
			}
		} else {
			// Add new worker
			worker.WeeklyProgressIdWorker = existingWeeklyProgress.ID
			if err := tx.Create(&worker).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Worker"})
				return
			}
		}
	}

	// Handle Materials: Update existing or add new materials
	for _, material := range weeklyProgress.Material {
		if material.ID != uuid.Nil {
			// Update existing material
			if err := tx.Model(&models.Material{}).Where("id = ?", material.ID).Updates(material).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Material"})
				return
			}
		} else {
			// Add new material
			material.WeeklyProgressIdMaterial = existingWeeklyProgress.ID
			if err := tx.Create(&material).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Material"})
				return
			}
		}
	}

	var materialIDsToDelete []uuid.UUID
	for _, existingMaterial := range existingWeeklyProgress.Material {
		found := false
		for _, updatedMaterial := range weeklyProgress.Material {
			if existingMaterial.ID == updatedMaterial.ID {
				found = true
				break
			}
		}
		if !found {
			// Mark for deletion
			materialIDsToDelete = append(materialIDsToDelete, existingMaterial.ID)
		}
	}

	// Delete marked materials
	if len(materialIDsToDelete) > 0 {
		if err := tx.Where("id IN ?", materialIDsToDelete).Delete(&models.Material{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Materials"})
			return
		}
	}

	// Step 2: Update existing materials and Step 3: Add new materials
	for _, material := range weeklyProgress.Material {
		if material.ID != uuid.Nil {
			// Update existing material
			if err := tx.Model(&models.Material{}).Where("id = ?", material.ID).Updates(material).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Material"})
				return
			}
		} else {
			// Add new material
			material.WeeklyProgressIdMaterial = existingWeeklyProgress.ID
			if err := tx.Create(&material).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Material"})
				return
			}
		}
	}

	// Commit the transaction

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"weekly_progress": existingWeeklyProgress,
		},
	})
}
