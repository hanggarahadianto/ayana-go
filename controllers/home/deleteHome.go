package controllers

import (
	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"

	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func ExtractPublicID(imageURL string) (string, error) {
// 	parsedURL, err := url.Parse(imageURL)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Hapus '/upload/vXYZ/' → ambil bagian setelahnya
// 	segments := strings.Split(parsedURL.Path, "/upload/")
// 	if len(segments) < 2 {
// 		return "", fmt.Errorf("format URL tidak valid")
// 	}

// 	publicID := segments[1]
// 	publicID = strings.TrimPrefix(publicID, "v")   // kadang masih ada v123/
// 	publicID = strings.SplitN(publicID, "/", 2)[1] // buang version
// 	publicID, _ = url.QueryUnescape(publicID)      // decode %20 → spasi

// 	// Hapus ekstensi ganda .png.png jika perlu
// 	publicID = strings.ReplaceAll(publicID, ".png.png", ".png")

// 	return publicID, nil
// }

// DeleteHome deletes a home from the database and removes its image from Cloudinary
func DeleteHome(c *gin.Context) {

	homeId := c.Param("id")

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

	publicID, err := uploadClaudinary.ExtractPublicID(home.Image)
	if err != nil {
		log.Println("❌ Gagal extract Public ID:", err)
	}

	err = uploadClaudinary.DeleteFromCloudinary(publicID)
	if err != nil {
		log.Println("❌ Gagal hapus dari Cloudinary:", err)
	}

	var infoIDs []uuid.UUID
	db.DB.Model(&models.Info{}).Where("home_id = ?", homeUUID).Pluck("id", &infoIDs)

	// Delete NearBy records linked to retrieved Info IDs
	if len(infoIDs) > 0 {
		fmt.Println("🗑 Deleting NearBy records linked to Info IDs:", infoIDs)
		db.DB.Debug().Where("info_id IN (?)", infoIDs).Delete(&models.NearBy{})
	}

	fmt.Println("🗑 Deleting Home record from database...")
	result := db.DB.Debug().Where("id = ?", homeUUID).Delete(&models.Home{})
	if result.Error != nil || result.RowsAffected == 0 {
		fmt.Println("❌ Failed to delete Home:", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Home ID doesn't exist or couldn't be deleted"})
		return
	}

	// fmt.Println("✅ Home and image deleted successfully")
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Home and image deleted successfully"})
}
