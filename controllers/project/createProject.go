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

	if err := c.ShouldBindJSON(&projectData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	newProject := models.Project{
		ProjectName:   projectData.ProjectName,
		ProjectLeader: projectData.ProjectLeader,
		Investor:      projectData.Investor,
		TotalCost:     projectData.TotalCost,
		ProjectTime:   projectData.ProjectTime,
		ProjectStart:  projectData.ProjectStart, // Set the ProjectStart value from the input
		ProjectEnd:    projectData.ProjectEnd,
		Note:          projectData.Note,
		CompanyID:     projectData.CompanyID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := db.DB.Create(&newProject)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newProject,
	})

}
