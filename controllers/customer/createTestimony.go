package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCustomerTestimony(c *gin.Context) {
	var input models.Testimony

	// 🔄 Binding JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// 🧠 Set ID jika kosong
	if input.ID == uuid.Nil {
		input.ID = uuid.New()
	}

	// 💾 Simpan ke DB
	// 💾 Simpan ke DB
	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan testimony"})
		return
	}

	// 🔄 Ambil ulang dengan Preload Customer
	var testimony models.Testimony
	if err := db.DB.Preload("Customer").First(&testimony, "id = ?", input.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil testimony setelah simpan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    testimony,
		"status":  "sukses",
		"message": "Testimony berhasil disimpan",
	})

}
