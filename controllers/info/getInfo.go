package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInfo(c *gin.Context) {

	var additionalInfo []models.Info

	infoId := c.Param("id")

	result := db.DB.Preload("NearBy").First(&additionalInfo, "home_id = ?", infoId)
	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		c.JSON(http.StatusOK, gin.H{
			"data":    additionalInfo,
			"status":  "error",
			"message": "Record not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   additionalInfo,
	})

}
