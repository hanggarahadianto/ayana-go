package controllers

import (
	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeleteHome deletes a home from the database and removes its image from Cloudinary
func DeleteHome(c *gin.Context) {
	fmt.Println("🚀 Starting DeleteHome process...")

	homeId := c.Param("id")
	fmt.Println("🔍 Home ID:", homeId)

	// Parse UUID
	homeUUID, err := uuid.Parse(homeId)
	if err != nil {
		fmt.Println("❌ Invalid Home ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Home ID"})
		return
	}

	// Fetch home details to get the image URL
	var home models.Home
	if err := db.DB.Debug().Where("id = ?", homeUUID).First(&home).Error; err != nil {
		fmt.Println("❌ Home not found in database:", err)
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Home not found"})
		return
	}

	// Log the image URL before deleting
	fmt.Println("🖼 Home Image URL:", home.Image)

	// Delete image from Cloudinary
	err = uploadClaudinary.DeleteFromCloudinary(home.Image)
	if err != nil {
		fmt.Println("❌ Error deleting image from Cloudinary:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete image from Cloudinary"})
		return
	}
	fmt.Println("✅ Image deleted from Cloudinary:", home.Image)

	// Delete related records
	fmt.Println("🗑 Deleting related Reservations and Info records...")
	db.DB.Debug().Where("home_id = ?", homeUUID).Delete(&models.Reservation{})
	db.DB.Debug().Where("home_id = ?", homeUUID).Delete(&models.Info{})

	// Get all Info IDs linked to this Home
	var infoIDs []uuid.UUID
	db.DB.Model(&models.Info{}).Where("home_id = ?", homeUUID).Pluck("id", &infoIDs)

	// Delete NearBy records linked to retrieved Info IDs
	if len(infoIDs) > 0 {
		fmt.Println("🗑 Deleting NearBy records linked to Info IDs:", infoIDs)
		db.DB.Debug().Where("info_id IN (?)", infoIDs).Delete(&models.NearBy{})
	}

	// Delete the Home itself
	fmt.Println("🗑 Deleting Home record from database...")
	result := db.DB.Debug().Where("id = ?", homeUUID).Delete(&models.Home{})
	if result.Error != nil || result.RowsAffected == 0 {
		fmt.Println("❌ Failed to delete Home:", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Home ID doesn't exist or couldn't be deleted"})
		return
	}

	fmt.Println("✅ Home and image deleted successfully")
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Home and image deleted successfully"})
}
