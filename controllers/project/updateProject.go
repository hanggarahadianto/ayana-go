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

	// Bind JSON body
	if err := c.ShouldBindJSON(&projectData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah project ada
	var existingProject models.Project
	if err := db.DB.First(&existingProject, projectData.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Project not found",
			"id":      projectData.ID,
		})
		return
	}

	// Update field dari request
	existingProject.Location = projectData.Location
	existingProject.Unit = projectData.Unit
	existingProject.Type = projectData.Type

	// Generate ulang project_name
	existingProject.ProjectName = projectData.Location + " - " + projectData.Unit + " - " + projectData.Type

	existingProject.ProjectLeader = projectData.ProjectLeader
	existingProject.Investor = projectData.Investor
	existingProject.TotalCost = projectData.TotalCost
	existingProject.ProjectTime = projectData.ProjectTime
	existingProject.ProjectStart = projectData.ProjectStart
	existingProject.ProjectEnd = projectData.ProjectEnd
	existingProject.Note = projectData.Note
	existingProject.ProjectStatus = projectData.ProjectStatus
	existingProject.CompanyID = projectData.CompanyID
	existingProject.UpdatedAt = time.Now()

	// Simpan ke database
	if err := db.DB.Save(&existingProject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existingProject,
	})
}
