package controllers

import (
	configClaudinary "ayana/config"
	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	utilsEnv "ayana/utils/env"
	"mime/multipart"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageUpdateResponse struct {
	ExistingImages []models.HomeImage `json:"existingImages"` // gambar yg tetap ada
	DeletedImages  []string           `json:"deletedImages"`  // ID gambar yg dihapus
	NewImages      []models.HomeImage `json:"newImages"`      // gambar baru yg berhasil diupload
}

func UpdateProductImages(c *gin.Context) {
	homeIdStr := c.Param("homeId")
	homeId, err := uuid.Parse(homeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "homeId tidak valid"})
		return
	}

	env, err := utilsEnv.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal load konfigurasi env"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca form"})
		return
	}

	keepImageIds := form.Value["keepImageIds"] // []string
	files := form.File["images"]               // []*multipart.FileHeader

	// Ambil semua gambar lama di DB
	var allImages []models.HomeImage
	if err := db.DB.Where("home_id = ?", homeId).Find(&allImages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data gambar"})
		return
	}

	// Track deleted images ID (yg gak ada di keepImageIds)
	var deletedImageIds []string
	if len(keepImageIds) == 0 {
		// hapus semua
		for _, img := range allImages {
			deletedImageIds = append(deletedImageIds, img.ID.String())
		}
		deleteAllImagesIfKeepEmpty(allImages, env)
	} else {
		for _, img := range allImages {
			if !contains(keepImageIds, img.ID.String()) {
				deletedImageIds = append(deletedImageIds, img.ID.String())
			}
		}
		deleteUnwantedImages(keepImageIds, allImages, env)
	}

	// Upload gambar baru dan simpan ke DB, track hasilnya
	var newImages []models.HomeImage
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		filePath := "products/" + homeId.String() + "/image_" + time.Now().Format("20060102150405")
		url, err := uploadClaudinary.UploadToCloudinary(file, filePath)
		if err != nil {
			continue
		}

		newImage := models.HomeImage{
			HomeID:    homeId,
			ImageURL:  url,
			CreatedAt: time.Now(),
		}
		if err := db.DB.Create(&newImage).Error; err == nil {
			newImages = append(newImages, newImage)
		}
	}

	// Ambil gambar yg masih ada di DB setelah update
	var existingImages []models.HomeImage
	if err := db.DB.Where("home_id = ?", homeId).Find(&existingImages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data gambar terbaru"})
		return
	}

	response := ImageUpdateResponse{
		ExistingImages: existingImages,
		DeletedImages:  deletedImageIds,
		NewImages:      newImages,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gambar berhasil diperbarui",
		"data":    response,
	})
}

func contains(list []string, id string) bool {
	for _, item := range list {
		if item == id {
			return true
		}
	}
	return false
}

func deleteAllImagesIfKeepEmpty(images []models.HomeImage, env utilsEnv.Config) {
	for _, image := range images {
		publicID := configClaudinary.GetPublicIDFromURL(
			image.ImageURL,
			configClaudinary.EnvCloudDeleteFolderHome(&env),
		)
		if publicID != "" {
			_ = uploadClaudinary.DeleteFromCloudinary(publicID)
		}
		db.DB.Delete(&image)
	}
}

// Fungsi untuk menghapus gambar yang tidak ada dalam keepImageIds
func deleteUnwantedImages(keepImageIds []string, allImages []models.HomeImage, env utilsEnv.Config) {
	for _, image := range allImages {
		if !contains(keepImageIds, image.ID.String()) {
			// Ambil publicID dari URL untuk dihapus dari Cloudinary
			publicID := configClaudinary.GetPublicIDFromURL(image.ImageURL, configClaudinary.EnvCloudDeleteFolderHome(&env))
			if publicID != "" {
				_ = uploadClaudinary.DeleteFromCloudinary(publicID)
			}
			// Hapus gambar dari DB
			db.DB.Delete(&image)
		}
	}
}

func uploadNewImages(files []*multipart.FileHeader, homeId uuid.UUID) {
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Upload path
		filePath := fmt.Sprintf("products/%s/image_%d", homeId.String(), time.Now().UnixNano())
		url, err := uploadClaudinary.UploadToCloudinary(file, filePath)
		if err != nil {
			continue
		}

		// Simpan gambar ke DB
		newImage := models.HomeImage{
			HomeID:    homeId,
			ImageURL:  url,
			CreatedAt: time.Now(),
		}
		_ = db.DB.Create(&newImage)
	}
}
