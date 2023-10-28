package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomes(c *gin.Context) {
	var homeList []models.Home

	result := db.DB.Debug().Find(&homeList)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   homeList,
	})

}
