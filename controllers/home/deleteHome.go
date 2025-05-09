package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteHome(c *gin.Context) {
	id := c.Param("id")

	var home models.Home
	if err := db.DB.Where("id = ?", id).First(&home).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Home not found"})
		return
	}

	// Optional: Hapus relasi NearBies terlebih dahulu jika diperlukan
	if err := db.DB.Where("home_id = ?", id).Delete(&models.NearBy{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related NearBies"})
		return
	}

	if err := db.DB.Delete(&home).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete home"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Home deleted successfully"})
}
