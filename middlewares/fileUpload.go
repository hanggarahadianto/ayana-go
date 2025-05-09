package controllers

import (
	"fmt"
	"net/http"
	"time"

	uploadClaudinary "ayana/utils/cloudinary-folder"

	"github.com/gin-gonic/gin"
)

// UploadImage handles image upload to Cloudinary
func UploadImage(c *gin.Context) {
	// Ambil file yang di-upload
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get the file",
		})
		return
	}

	// Tentukan path file di Cloudinary, bisa menggunakan UUID atau nama file yang unik
	filePath := "uploads/" + "image_" + fmt.Sprintf("%d", time.Now().UnixNano()) // contoh penamaan berdasarkan timestamp

	// Upload file ke Cloudinary
	uploadedURL, err := uploadClaudinary.UploadToCloudinary(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image to Cloudinary",
		})
		return
	}

	// Jika berhasil, return URL gambar yang sudah di-upload
	c.JSON(http.StatusOK, gin.H{
		"message":   "Image uploaded successfully",
		"image_url": uploadedURL,
	})
}
