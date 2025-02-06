package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteWeeklyProgress(c *gin.Context) {
	// Get the project ID from the URL parameter
	weeklyProgressID := c.Param("id")

	// Start a transaction
	tx := db.DB.Begin()

	// Ensure the transaction is rolled back if an error occurs
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Internal server error",
			})
		}
	}()

	// Find WeeklyProgress
	var weeklyProgress models.WeeklyProgress
	if err := tx.Where("id = ?", weeklyProgressID).First(&weeklyProgress).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Weekly Progress not found",
		})
		return
	}

	// Log the found WeeklyProgress for debugging
	log.Printf("Weekly Progress found: %+v", weeklyProgress)

	// Check if related materials exist
	var materials []models.Material
	if err := tx.Where("weekly_progress_id_material = ?", weeklyProgressID).Find(&materials).Error; err != nil {
		tx.Rollback()
		log.Printf("Error fetching materials: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Error fetching related materials",
		})
		return
	}
	log.Printf("Materials found: %+v", materials)

	// Delete related records (Material)
	if err := tx.Unscoped().Where("weekly_progress_id_material = ?", weeklyProgressID).Delete(&models.Material{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting material: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to delete related materials",
		})
		return
	}

	var workers []models.Worker
	if err := tx.Where("weekly_progress_id_worker = ?", weeklyProgressID).Find(&workers).Error; err != nil {
		tx.Rollback()
		log.Printf("Error fetching workers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Error fetching related workers",
		})
		return
	}
	log.Printf("Workers found: %+v", materials)

	// Delete related records (Worker)
	if err := tx.Unscoped().Where("weekly_progress_id_worker = ?", weeklyProgressID).Delete(&models.Worker{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting worker: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to delete related workers",
		})
		return
	}

	// Delete the WeeklyProgress record
	if err := tx.Delete(&weeklyProgress).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting weekly progress: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to delete weekly progress",
		})
		return
	}

	// Commit the transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Weekly Progress deleted successfully",
	})
}
