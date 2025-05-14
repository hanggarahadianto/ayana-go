package controllers

import (
	"ayana/db"
	"ayana/dto" // ✅ pakai dto
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeById(c *gin.Context) {
	homeId := c.Param("id")

	var home models.Home
	if err := db.DB.Preload("NearBies").Preload("Cluster").First(&home, "id = ?", homeId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Home id doesn't exist",
		})
		return
	}

	// Map model ke DTO
	response := dto.HomeByClusterResponse{
		ID:         home.ID.String(),
		Title:      home.Title,
		Type:       home.Type,
		Maps:       home.Cluster.Maps,
		Content:    home.Content,
		Price:      home.Price,
		Status:     home.Status,
		Square:     home.Square,
		Bedroom:    home.Bedroom,
		Bathroom:   home.Bathroom,
		StartPrice: home.StartPrice,
		Cluster: dto.ClusterResponse{
			ID:       home.Cluster.ID.String(),
			Name:     home.Cluster.Name,
			Location: home.Cluster.Location,
			Maps:     home.Cluster.Maps,
		},
	}

	// Map nearbies
	for _, nearby := range home.NearBies {
		response.NearBies = append(response.NearBies, dto.NearBy{
			ID:       nearby.ID.String(),
			Name:     nearby.Name,
			Distance: nearby.Distance,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}
