package controllers

import (
	"ayana/db"
	"ayana/models"
	cloudinaryConfig "ayana/utils/cloudinary-folder"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func DeleteHome(c *gin.Context) {

// 	homeId := c.Param("id")

// 	homeUUID, err := uuid.Parse(homeId)
// 	if err != nil {
// 		fmt.Println("❌ Invalid Home ID:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Home ID"})
// 		return
// 	}

// 	// Fetch home details to get the image URL
// 	var home models.Home
// 	if err := db.DB.Debug().Where("id = ?", homeUUID).First(&home).Error; err != nil {
// 		fmt.Println("❌ Home not found in database:", err)
// 		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Home not found"})
// 		return
// 	}

// 	fmt.Println("✅ Found Home:", home)

// 	publicID, err := uploadClaudinary.ExtractPublicID(home.Image)
// 	if err != nil {
// 		log.Println("❌ Gagal extract Public ID:", err)
// 	}

// 	err = uploadClaudinary.DeleteFromCloudinary(publicID)
// 	if err != nil {
// 		log.Println("❌ Gagal hapus dari Cloudinary:", err)
// 	}

// 	var infoIDs []uuid.UUID
// 	db.DB.Model(&models.Info{}).Where("home_id = ?", homeUUID).Pluck("id", &infoIDs)

// 	if len(infoIDs) > 0 {
// 		fmt.Println("🗑 Deleting NearBy records linked to Info IDs:", infoIDs)
// 		db.DB.Debug().Where("info_id IN (?)", infoIDs).Delete(&models.NearBy{})
// 	}

// 	// Delete reservation(s) before deleting the home
// 	fmt.Println("🗑 Deleting Reservation records linked to Home ID:", homeUUID)
// 	if result := db.DB.Debug().Where("home_id = ?", homeUUID).Delete(&models.Reservation{}); result.Error != nil {
// 		fmt.Println("❌ Failed to delete reservations:", result.Error)
// 	}

// 	// Delete the home record itself
// 	fmt.Println("🗑 Deleting Home record from database...")
// 	result := db.DB.Debug().Unscoped().Where("id = ?", homeUUID).Delete(&models.Home{})
// 	if result.Error != nil || result.RowsAffected == 0 {
// 		fmt.Println("❌ Failed to delete Home:", result.Error)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Home ID doesn't exist or couldn't be deleted"})
// 		return
// 	}
// 	// fmt.Println("✅ Home and image deleted successfully")
// 	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Home and image deleted successfully"})
// }

func DeleteHome(c *gin.Context) {
	homeID := c.Param("id")

	// Validasi UUID
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		log.Printf("❌ Invalid Home ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "ID Home tidak valid"})
		return
	}

	// Mulai transaksi database
	tx := db.DB.Begin()
	if tx.Error != nil {
		log.Printf("❌ Gagal memulai transaksi: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal memproses permintaan"})
		return
	}
	defer tx.Rollback()

	// Cari data Home
	var home models.Home
	if err := tx.Where("id = ?", homeUUID).First(&home).Error; err != nil {
		log.Printf("❌ Home tidak ditemukan: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Home tidak ditemukan"})
		return
	}

	log.Printf("✅ Ditemukan Home: %+v", home)

	// Hapus gambar dari Cloudinary jika ada
	if home.Image != "" {
		if err := cloudinaryConfig.DeleteFromCloudinary(home.Image); err != nil {
			log.Printf("⚠️ Gagal menghapus gambar dari Cloudinary: %v", err)
			// Lanjutkan penghapusan database meskipun gagal hapus gambar
		} else {
			log.Println("✅ Gambar berhasil dihapus dari Cloudinary")
		}
	} else {
		log.Println("⚠️ Tidak ada gambar untuk dihapus")
	}

	// Hapus data dari database
	if err := tx.Delete(&home).Error; err != nil {
		log.Printf("❌ Gagal menghapus Home dari database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal menghapus data"})
		return
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		log.Printf("❌ Gagal commit transaksi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal menyelesaikan penghapusan"})
		return
	}

	log.Println("✅ Home berhasil dihapus dari database")
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Home berhasil dihapus"})
}
