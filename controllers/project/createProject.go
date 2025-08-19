package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateProject(c *gin.Context) {
	var projectData models.Project

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&projectData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	// Generate project_name dari location + unit + type
	projectName := projectData.Location + " - " + projectData.Unit + " - " + projectData.Type

	newProject := models.Project{
		ProjectName:   projectName,
		Location:      projectData.Location,
		Type:          projectData.Type,
		Unit:          projectData.Unit,
		ProjectLeader: projectData.ProjectLeader,
		Investor:      projectData.Investor,
		ProjectTime:   projectData.ProjectTime,
		TotalCost:     projectData.TotalCost,
		ProjectStart:  projectData.ProjectStart,
		ProjectEnd:    projectData.ProjectEnd,
		ProjectStatus: projectData.ProjectStatus,
		Note:          projectData.Note,
		CompanyID:     projectData.CompanyID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Simpan ke database
	if err := db.DB.Create(&newProject).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newProject,
	})
}
