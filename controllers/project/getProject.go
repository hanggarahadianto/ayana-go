package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProject(c *gin.Context) {
	var projectList []models.Project

	result := db.DB.Debug().Order("created_at desc, updated_at desc").Find(&projectList)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projectList,
	})

}
