package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateHome(c *gin.Context) {
	var input models.Home
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	home := models.Home{
		ID:         uuid.New(),
		Title:      input.Title,
		Type:       input.Type,
		Content:    input.Content,
		Bathroom:   input.Bathroom,
		Bedroom:    input.Bedroom,
		Square:     input.Square,
		Price:      input.Price,
		Quantity:   input.Quantity,
		Status:     input.Status,
		Sequence:   input.Sequence,
		StartPrice: input.StartPrice,
		ClusterID:  input.ClusterID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	for i := range input.NearBies {
		input.NearBies[i].ID = uuid.New()
		input.NearBies[i].HomeID = home.ID
	}

	home.NearBies = input.NearBies

	if err := db.DB.Create(&home).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create home"})
		return
	}
	c.JSON(http.StatusCreated, home)
}
