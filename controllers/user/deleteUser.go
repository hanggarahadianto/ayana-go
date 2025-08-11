package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Deleteuser menghapus akun dan relasi TransactionCategory terkait
func DeleteUser(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	username, _ := c.Get("username")
	if username != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Hanya superadmin yang dapat mengakses data ini",
			"status":  "error",
		})
		return
	}

	// Periksa apakah akun ada
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Hapus akun itu sendiri
	if err := db.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user  deleted successfully",
	})
}
