package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateCustomerTestimony(c *gin.Context) {
	var input models.Testimony
	id := c.Param("id")

	// 🔄 Binding JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Pastikan testimony ada di DB
	var existingTestimony models.Testimony
	if err := db.DB.First(&existingTestimony, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Testimony tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil testimony"})
		return
	}

	// ✅ Cek apakah customer_id valid
	var customer models.Customer
	if err := db.DB.First(&customer, "id = ?", input.CustomerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data customer"})
		return
	}

	// ✏️ Update field (langsung update data pada objek existing)
	existingTestimony.CustomerID = input.CustomerID
	existingTestimony.Rating = input.Rating
	existingTestimony.Note = input.Note

	// 💾 Simpan perubahan
	if err := db.DB.Save(&existingTestimony).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui testimony"})
		return
	}

	// 🔄 Ambil ulang dengan Preload Customer
	var updatedTestimony models.Testimony
	if err := db.DB.Preload("Customer").First(&updatedTestimony, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil testimony setelah update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedTestimony,
		"status":  "sukses",
		"message": "Testimony berhasil diperbarui",
	})
}
