package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateHome(c *gin.Context) {
	id := c.Param("id")
	var home models.Home

	// Cari data berdasarkan ID
	if err := db.DB.First(&home, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Home not found",
		})
		return
	}

	// // Ambil nilai dari form
	// title := c.Request.PostFormValue("title")
	// location := c.Request.PostFormValue("location")
	// content := c.Request.PostFormValue("content")
	// address := c.Request.PostFormValue("address")
	// bathroom := c.Request.PostFormValue("bathroom")
	// bedroom := c.Request.PostFormValue("bedroom")
	// square := c.Request.PostFormValue("square")
	// priceStr := c.Request.PostFormValue("price")
	// quantityStr := c.Request.PostFormValue("quantity")
	// status := c.Request.PostFormValue("status")
	// sequenceStr := c.Request.PostFormValue("sequence")

	// // Validasi dan konversi numerik
	// price, err := strconv.Atoi(priceStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid price"})
	// 	return
	// }
	// quantity, err := strconv.Atoi(quantityStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid quantity"})
	// 	return
	// }
	// sequence, err := strconv.Atoi(sequenceStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid sequence"})
	// 	return
	// }

	// // Cek apakah file baru diunggah
	// file, header, err := c.Request.FormFile("file")
	// var imageUrl string
	// if err == nil {
	// 	defer file.Close()

	// 	// Upload file baru
	// 	imageUrl, err = uploadClaudinary.UploadToCloudinary(file, header.Filename)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to upload image"})
	// 		return
	// 	}

	// 	// Hapus gambar lama
	// 	if home.Image != "" {
	// 		_ = uploadClaudinary.DeleteFromCloudinary(home.Image)
	// 	}
	// } else {
	// 	// Jika tidak upload file, gunakan gambar lama
	// 	imageUrl = home.Image
	// }

	// // Update data
	// home.Title = title
	// home.Location = location
	// home.Content = content
	// home.Address = address
	// home.Bathroom = bathroom
	// home.Bedroom = bedroom
	// home.Square = square
	// home.Price = float64(price)
	// home.Quantity = quantity
	// home.Status = status
	// home.Sequence = sequence
	// home.Image = imageUrl
	// home.UpdatedAt = time.Now()

	// if err := db.DB.Save(&home).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to update home"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Home updated successfully",
		"data":    home,
	})
}
