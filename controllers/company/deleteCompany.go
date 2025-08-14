package controllers

import (
	"ayana/db"
	"ayana/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	username, _ := c.Get("username")
	if username != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Hanya superadmin yang dapat mengakses data ini",
			"status":  "error",
		})
		return
	}

	// Periksa apakah company ada
	var company models.Company
	if err := db.DB.First(&company, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	// Hapus dari database
	if err := db.DB.Delete(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company deleted successfully"})
}
