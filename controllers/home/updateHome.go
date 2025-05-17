package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateHome(c *gin.Context) {
	var input models.Home
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak ditemukan dalam payload"})
		return
	}

	var home models.Home
	if err := db.DB.First(&home, "id = ?", input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Home tidak ditemukan"})
		return
	}

	// Perbarui field
	home.Title = input.Title
	home.Content = input.Content
	home.Type = input.Type
	home.Bathroom = input.Bathroom
	home.Bedroom = input.Bedroom
	home.Square = input.Square
	home.Price = input.Price
	home.Quantity = input.Quantity
	home.Status = input.Status
	home.Sequence = input.Sequence
	home.StartPrice = input.StartPrice
	home.ClusterID = input.ClusterID
	home.UpdatedAt = time.Now()

	// Simpan data home yang sudah diperbarui
	if err := db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&home).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui home"})
		return
	}

	c.JSON(http.StatusOK, home)
}
