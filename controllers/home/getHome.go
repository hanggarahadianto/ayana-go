package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomes(c *gin.Context) {
	var homeList []models.Home

	result := db.DB.Order("created_at desc, updated_at desc").Find(&homeList)
	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
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
