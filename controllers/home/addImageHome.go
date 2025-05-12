package controllers

import (
	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadProductImage handles image upload to Cloudinary and saves image info to the database
func UploadProductImage(c *gin.Context) {
	homeIdStr := c.Param("homeId")

	// Konversi string ke uuid.UUID
	homeId, err := uuid.Parse(homeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "homeId tidak valid"})
		return
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca multipart form"})
		return
	}

	files := form.File["images"] // nama field harus "images"

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada file yang diupload"})
		return
	}

	var uploadedURLs []string

	for _, fileHeader := range files {
		homeImage := models.HomeImage{
			HomeID:    homeId,
			CreatedAt: time.Now(),
		}

		if err := db.DB.Create(&homeImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal menyimpan data gambar ke database: %v", err)})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		filePath := fmt.Sprintf("products/%s/image_%d", homeId.String(), time.Now().UnixNano())
		url, err := uploadClaudinary.UploadToCloudinary(file, filePath)
		if err != nil {
			db.DB.Delete(&homeImage)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal upload gambar ke Cloudinary: %v", err)})
			return
		}

		homeImage.ImageURL = url

		if err := db.DB.Save(&homeImage).Error; err != nil {
			db.DB.Delete(&homeImage)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal mengupdate URL gambar di database: %v", err)})
			return
		}

		uploadedURLs = append(uploadedURLs, url)
	}

	if len(uploadedURLs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Semua gambar gagal di-upload ke Cloudinary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Berhasil upload beberapa gambar",
		"image_urls": uploadedURLs,
	})
}
