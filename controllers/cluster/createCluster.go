package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCluster(c *gin.Context) {
	var input models.Cluster

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cluster := models.Cluster{
		ID:        uuid.New(),
		Name:      input.Name,
		Location:  input.Location,
		Square:    input.Square,
		Price:     input.Price,
		Quantity:  input.Quantity,
		Status:    input.Status,
		Sequence:  input.Sequence,
		Maps:      input.Maps,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cluster"})
		return
	}

	c.JSON(http.StatusCreated, cluster)
}
