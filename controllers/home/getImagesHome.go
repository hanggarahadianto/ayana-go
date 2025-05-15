package controllers

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomeImages(c *gin.Context) {
	homeId := c.Param("homeId") // Correctly retrieve homeId from the URL

	// Log for debugging purposes
	fmt.Println("homeId:", homeId)

	var home models.Home // Assuming Home is a model, not a slice

	// Retrieve home record by ID
	result := db.DB.First(&home, "id = ?", homeId)
	if result.Error != nil {
		// Return error if homeId is not found
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": fmt.Sprintf("Home with ID %s doesn't exist", homeId),
		})
		return
	}

	// Retrieve related HomeImages using a HasMany relationship (assuming it's set in the Home model)
	var homeImages []models.HomeImage
	db.DB.Where("home_id = ?", homeId).Find(&homeImages)

	// Set the first image URL as the thumbnail if available
	var thumbnail string
	if len(homeImages) > 0 {
		thumbnail = homeImages[0].ImageURL
	}

	// Build array of objects with id + url
	var imageData []gin.H
	for _, image := range homeImages {
		imageData = append(imageData, gin.H{
			"id":  image.ID,
			"url": image.ImageURL,
		})
	}

	// Return the home images and the thumbnail
	c.JSON(http.StatusOK, gin.H{
		"images":    imageData,
		"thumbnail": thumbnail,
	})

}
