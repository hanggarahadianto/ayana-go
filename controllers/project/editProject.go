package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func EditProject(c *gin.Context) {
	var projectData models.Project

	// Bind JSON body to projectData
	if err := c.ShouldBindJSON(&projectData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the project exists using the ID from the body
	var existingProject models.Project
	if err := db.DB.First(&existingProject, projectData.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Project not found", "id": projectData.ID})
		return
	}

	// Update project data with new values from the request body
	existingProject.ProjectName = projectData.ProjectName
	existingProject.ProjectLeader = projectData.ProjectLeader
	existingProject.Investor = projectData.Investor
	existingProject.TotalCost = projectData.TotalCost
	existingProject.ProjectTime = projectData.ProjectTime
	existingProject.ProjectStart = projectData.ProjectStart
	existingProject.ProjectEnd = projectData.ProjectEnd
	existingProject.Note = projectData.Note
	existingProject.UpdatedAt = time.Now()

	// Save the updated project to the database
	if err := db.DB.Save(&existingProject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existingProject,
	})
}
