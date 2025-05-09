package controllers

import (
	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadProductImage handles image upload to Cloudinary and saves image info to the database
func UploadProductImage(c *gin.Context) {
	homeId := c.Param("homeId") // Perbaiki dari productId menjadi homeId

	// Menambahkan log untuk homeId
	fmt.Println("homeId:", homeId)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca multipart form"})
		return
	}

	files := form.File["images"] // nama field harus "images" di Postman atau FE

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada file yang diupload"})
		return
	}

	var uploadedURLs []string

	// Pastikan semua gambar disimpan ke database terlebih dahulu
	for _, fileHeader := range files {

		homeImage := models.HomeImage{
			HomeID:    homeId,
			CreatedAt: time.Now(),
		}

		// Simpan gambar ke database terlebih dahulu
		if err := db.DB.Create(&homeImage).Error; err != nil {
			// Jika gagal menyimpan data gambar ke database, kirim error dan hentikan proses
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal menyimpan data gambar ke database: %v", err)})
			return
		}

		// Upload gambar ke Cloudinary hanya jika berhasil simpan ke database
		file, err := fileHeader.Open()
		if err != nil {
			// Jika gagal membuka file, lanjutkan ke file berikutnya
			continue
		}
		defer file.Close()

		// Tentukan path untuk file gambar di Cloudinary
		filePath := fmt.Sprintf("products/%s/image_%d", homeId, time.Now().UnixNano())
		url, err := uploadClaudinary.UploadToCloudinary(file, filePath)
		if err != nil {
			// Jika gagal upload ke Cloudinary, rollback entri gambar dari database
			db.DB.Delete(&homeImage)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal upload gambar ke Cloudinary: %v", err)})
			return
		}

		// Update URL gambar setelah berhasil di-upload ke Cloudinary
		homeImage.ImageURL = url

		// Update data gambar dengan URL Cloudinary ke database
		if err := db.DB.Save(&homeImage).Error; err != nil {
			// Jika gagal mengupdate URL gambar, rollback entri gambar dari database
			db.DB.Delete(&homeImage)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal mengupdate URL gambar di database: %v", err)})
			return
		}

		// Append URL gambar yang berhasil di-upload ke Cloudinary
		uploadedURLs = append(uploadedURLs, url)
	}

	// Jika tidak ada gambar yang berhasil di-upload ke Cloudinary, kirimkan error
	if len(uploadedURLs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Semua gambar gagal di-upload ke Cloudinary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Berhasil upload beberapa gambar",
		"image_urls": uploadedURLs,
	})
}
