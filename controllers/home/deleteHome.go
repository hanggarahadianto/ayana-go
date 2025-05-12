package controllers

import (
	configClaudinary "ayana/config"
	"ayana/db"
	"ayana/models"

	uploadClaudinary "ayana/utils/cloudinary-folder"

	utilsEnv "ayana/utils/env"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteHome deletes a home and all its images (DB + Cloudinary)
func DeleteHome(c *gin.Context) {
	homeId := c.Param("homeId")

	// Step 1: Cek apakah Home ada
	var home models.Home
	if err := db.DB.First(&home, "id = ?", homeId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Home dengan ID %s tidak ditemukan", homeId)})
		return
	}

	// Step 2: Ambil semua gambar terkait dari DB
	var homeImages []models.HomeImage
	if err := db.DB.Where("home_id = ?", homeId).Find(&homeImages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data gambar"})
		return
	}

	// Step 3: Hapus gambar dari Cloudinary
	env, err := utilsEnv.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal load environment config"})
		return
	}

	for _, image := range homeImages {
		publicID := configClaudinary.GetPublicIDFromURL(image.ImageURL, configClaudinary.EnvCloudDeleteFolderHome(&env))
		if publicID == "" {
			log.Printf("❌ PublicID tidak ditemukan untuk URL: %s", image.ImageURL)
			continue
		}

		if err := uploadClaudinary.DeleteFromCloudinary(publicID); err != nil {
			log.Printf("❌ Gagal hapus image dari Cloudinary: %v", err)
			// Lanjutkan menghapus sisanya walaupun error
		}
	}

	// Step 4: Hapus record gambar dari DB
	if err := db.DB.Where("home_id = ?", homeId).Delete(&models.HomeImage{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data gambar di database"})
		return
	}

	// Step 5: Hapus Home dari DB
	if err := db.DB.Delete(&home).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus home dari database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Home dan semua gambar berhasil dihapus"})
}
