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
	"gorm.io/gorm"
)

func DeleteCluster(c *gin.Context) {
	clusterId := c.Param("id")

	// Load env config
	env, err := utilsEnv.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal load environment config"})
		return
	}

	// Start transaction
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Ambil cluster + preload semua relasi homes + images
		var cluster models.Cluster
		if err := tx.Preload("Homes.Images").First(&cluster, "id = ?", clusterId).Error; err != nil {

			return fmt.Errorf("cluster does not exist")
		}

		// Step 2: Hapus seluruh home + gambar dari DB & Cloudinary
		for _, home := range cluster.Homes {
			for _, image := range home.Images {
				publicID := config.GetPublicIDFromURL(image.ImageURL, config.EnvCloudDeleteFolderHome(&env))
				if publicID != "" {
					if err := uploadClaudinary.DeleteFromCloudinary(publicID); err != nil {
						log.Printf("❌ Gagal hapus image dari Cloudinary (Home ID: %s): %v", home.ID, err)
					}
				}
			}

			// Hapus semua HomeImage dari DB
			if err := tx.Where("home_id = ?", home.ID).Delete(&models.HomeImage{}).Error; err != nil {
				return fmt.Errorf("❌ Gagal hapus HomeImages: %v", err)
			}

			// Hapus Home
			if err := tx.Unscoped().Delete(&home).Error; err != nil {
				return fmt.Errorf("❌ Gagal hapus Home: %v", err)
			}
		}

		// Step 3: Hapus Cluster (hard delete)
		if err := tx.Unscoped().Delete(&cluster).Error; err != nil {
			return fmt.Errorf("❌ Gagal hapus Cluster: %v", err)
		}

		return nil
	})

	// Jika ada error dalam transaction
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Berhasil
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Cluster %s dan seluruh data terkait berhasil dihapus", clusterId)})
}
