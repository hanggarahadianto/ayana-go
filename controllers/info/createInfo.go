package controllers

import (
	"net/http"

	"ayana/db"
	"ayana/models"

	"time"

	"github.com/gin-gonic/gin"
)

func CreateInfo(c *gin.Context) {

	var infoData models.Info

	if err := c.ShouldBindJSON(&infoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingInfo models.Info
	if err := db.DB.Where("home_id = ?", infoData.HomeID).First(&existingInfo).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Info record with this home_id already exists",
		})
		return
	}

	now := time.Now()

	newInfo := models.Info{
		Maps:       infoData.Maps,
		StartPrice: infoData.StartPrice,
		HomeID:     infoData.HomeID,
		NearBy:     infoData.NearBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	db.DB.Exec("DISCARD ALL")

	result := db.DB.Debug().Create(&newInfo)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newInfo,
	})

}
