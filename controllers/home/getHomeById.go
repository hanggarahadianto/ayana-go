package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeById(c *gin.Context) {
	homeId := c.Param("id")

	var home models.Home

	result := db.DB.Debug().First(&home, "id = ?", homeId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Home id doesn't exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   home,
	})
}
