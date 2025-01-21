package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProjectById(c *gin.Context) {
	projectId := c.Param("id")

	var project models.Project

	result := db.DB.Debug().First(&project, "id = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Home id doesn't exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   project,
	})
}
