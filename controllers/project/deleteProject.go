package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteProject(c *gin.Context) {
	// Get the project ID from the URL parameter
	projectID := c.Param("id")

	// Find the project by ID
	var project models.Project
	if err := db.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Project not found",
		})
		return
	}

	// Delete the project
	if err := db.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to delete project",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project deleted successfully",
	})
}
