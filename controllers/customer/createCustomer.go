package controllers

import (
	"ayana/db"
	"ayana/models"
	customer "ayana/service/customer"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCustomer(c *gin.Context) {
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Ambil data marketer dari DB berdasarkan MarketerID
	var marketer models.Employee
	if err := db.DB.First(&marketer, "id = ?", input.MarketerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Marketer tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data marketer"})
		return
	}

	// ✅ Isi marketer_name dari relasi
	input.MarketerName = marketer.Name

	// 🔄 Panggil service untuk buat customer
	createdCustomer, err := customer.CreateCustomer(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat customer: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCustomer)
}
