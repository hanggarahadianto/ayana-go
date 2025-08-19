package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// FinishProject menandai project sebagai selesai
type FinishProjectInput struct {
	ProjectFinished string `json:"project_finished" binding:"required"`
	ProjectStatus   string `json:"project_status" binding:"required"`
}

func FinishProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := db.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Project not found"})
		return
	}

	var input FinishProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// parse YYYY-MM-DD → time.Time
	t, err := time.Parse("2006-01-02", input.ProjectFinished)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, must be YYYY-MM-DD"})
		return
	}

	project.ProjectFinished = &t
	project.ProjectStatus = models.StatusDone
	project.UpdatedAt = time.Now()

	if err := db.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project marked as finished",
		"data":    project,
	})
}
