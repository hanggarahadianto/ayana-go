package controllers

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeById(c *gin.Context) {
	homeId := c.Param("id")

	var home models.Home
	if err := db.DB.
		Preload("Cluster").
		Preload("Cluster.NearBies").
		First(&home, "id = ?", homeId).Error; err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Home id doesn't exist",
		})
		return
	}

	// Validasi jika Cluster tidak ada
	if home.Cluster == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Cluster for this home not found",
		})
		return
	}

	// Map NearBies ke DTO
	var nearBiesDTO []dto.NearBy
	for _, nearby := range home.Cluster.NearBies {
		nearBiesDTO = append(nearBiesDTO, dto.NearBy{
			ID:       nearby.ID.String(),
			Name:     nearby.Name,
			Distance: nearby.Distance,
		})
	}

	// Map Home + Cluster + NearBies ke response DTO
	response := dto.HomeByClusterResponse{
		ID:         home.ID.String(),
		Title:      home.Title,
		Type:       home.Type,
		Maps:       home.Cluster.Maps,
		Content:    home.Content,
		Price:      home.Price,
		Location:   home.Cluster.Location,
		Status:     home.Status,
		Square:     home.Square,
		Bedroom:    home.Bedroom,
		Bathroom:   home.Bathroom,
		StartPrice: home.StartPrice,
		Quantity:   home.Quantity,
		Sequence:   home.Sequence,
		Cluster: dto.ClusterResponse{
			ID:   home.Cluster.ID.String(),
			Name: home.Cluster.Name,
			Maps: home.Cluster.Maps,
		},
		NearBies: nearBiesDTO,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}
