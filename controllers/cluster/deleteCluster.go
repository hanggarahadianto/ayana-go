package controllers

import (
	"ayana/config"
	"ayana/db"
	"ayana/models"

	uploadClaudinary "ayana/utils/cloudinary-folder"
	utilsEnv "ayana/utils/env"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCluster(c *gin.Context) {
	clusterId := c.Param("id")

	// Step 1: Cari cluster dan preload homes
	var cluster models.Cluster
	if err := db.DB.Preload("Homes.HomeImages").First(&cluster, "id = ?", clusterId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Step 2: Load env
	env, err := utilsEnv.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal load environment config"})
		return
	}

	// Step 3: Hapus semua homes dan image-nya
	for _, home := range cluster.Homes {
		// Hapus gambar dari Cloudinary
		for _, image := range home.Images {
			publicID := config.GetPublicIDFromURL(image.ImageURL, config.EnvCloudDeleteFolderHome(&env))
			if publicID != "" {
				if err := uploadClaudinary.DeleteFromCloudinary(publicID); err != nil {
					log.Printf("❌ Gagal hapus image dari Cloudinary: %v", err)
				}
			} else {
				log.Printf("❌ PublicID tidak ditemukan untuk URL: %s", image.ImageURL)
			}
		}

		// Hapus record gambar dari DB
		if err := db.DB.Where("home_id = ?", home.ID).Delete(&models.HomeImage{}).Error; err != nil {
			log.Printf("❌ Gagal menghapus gambar dari DB untuk home %s: %v", home.ID, err)
		}

		// Hapus home
		if err := db.DB.Delete(&home).Error; err != nil {
			log.Printf("❌ Gagal menghapus home %s: %v", home.ID, err)
		}
	}

	// Step 4: Hapus cluster
	if err := db.DB.Delete(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cluster"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Cluster %s dan seluruh data terkait berhasil dihapus", clusterId)})
}
