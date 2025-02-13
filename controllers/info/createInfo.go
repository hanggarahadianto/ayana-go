package controllers

import (
	"net/http"

	"ayana/db"
	"ayana/models"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Create the Info record first
	result := db.DB.Create(&newInfo)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	// Explicitly set InfoID for each NearBy record and insert
	var nearByRecords []models.NearBy
	for _, near := range infoData.NearBy {
		nearByRecords = append(nearByRecords, models.NearBy{
			ID:       uuid.New(), // Generate UUID for NearBy
			Name:     near.Name,
			Distance: near.Distance,
			InfoID:   newInfo.ID, // Set foreign key
		})
	}

	// Insert all NearBy records in batch
	if len(nearByRecords) > 0 {
		if err := db.DB.Create(&nearByRecords).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Failed to create NearBy records: " + err.Error(),
			})
			return
		}
	}

	// Fetch the newly created record with associations
	db.DB.Preload("NearBy").First(&newInfo, newInfo.ID)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newInfo,
	})
}
