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

	if err := c.ShouldBindJSON(&weeklyProgress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project models.Project
	if err := db.DB.Where("id = ?", weeklyProgress.ProjectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking Project ID"})
		}
		return
	}

	var existingWeeklyProgress models.WeeklyProgress
	if err := db.DB.Where("id = ?", weeklyProgress.ID).Preload("Material").Preload("Worker").First(&existingWeeklyProgress).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Weekly Progress ID does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Weekly Progress"})
		}
		return
	}

	tx := db.DB.Begin()

	existingWeeklyProgress.WeekNumber = weeklyProgress.WeekNumber
	existingWeeklyProgress.Percentage = weeklyProgress.Percentage
	existingWeeklyProgress.AmountMaterial = weeklyProgress.AmountMaterial
	existingWeeklyProgress.AmountWorker = weeklyProgress.AmountWorker
	existingWeeklyProgress.Note = weeklyProgress.Note
	existingWeeklyProgress.ProjectID = weeklyProgress.ProjectID
	existingWeeklyProgress.UpdatedAt = time.Now()

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

	for _, worker := range weeklyProgress.Worker {
		if worker.ID != uuid.Nil {
			if err := tx.Model(&models.Worker{}).Where("id = ?", worker.ID).Updates(worker).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Worker"})
				return
			}
		} else {
			worker.WeeklyProgressIdWorker = existingWeeklyProgress.ID
			if err := tx.Create(&worker).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Worker"})
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
			materialIDsToDelete = append(materialIDsToDelete, existingMaterial.ID)
		}
	}

	if len(materialIDsToDelete) > 0 {
		if err := tx.Where("id IN ?", materialIDsToDelete).Delete(&models.Material{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Materials"})
			return
		}
	}

	// ✅ Only one loop for materials (Fixed)
	for _, material := range weeklyProgress.Material {
		if material.ID != uuid.Nil {
			if err := tx.Model(&models.Material{}).Where("id = ?", material.ID).Updates(material).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Material"})
				return
			}
		} else {
			material.WeeklyProgressIdMaterial = existingWeeklyProgress.ID
			if err := tx.Create(&material).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Material"})
				return
			}
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"weekly_progress": existingWeeklyProgress,
		},
	})
}
