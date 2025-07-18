package controllers

// import (
// 	"ayana/db"
// 	"ayana/models"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// func UpdateTestimony(c *gin.Context) {
// 	id := c.Param("id")

// 	// ✅ Parse UUID
// 	testimonyID, err := uuid.Parse(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	// 🔍 Ambil data testimony yang mau diupdate
// 	var existing models.Testimony
// 	if err := db.DB.First(&existing, "id = ?", testimonyID).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Testimony tidak ditemukan"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data testimony"})
// 		return
// 	}

// 	// 🔄 Bind input yang akan diupdate
// 	var input models.Testimony
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// ✅ Validasi rating
// 	if input.Rating < 1 || input.Rating > 5 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating harus antara 1 sampai 5"})
// 		return
// 	}

// 	// ✅ Optional: validasi HomeID kalau diubah
// 	if input.HomeID != nil {
// 		var dummy struct{}
// 		if err := db.DB.Table("homes").First(&dummy, "id = ?", input.HomeID).Error; err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": "HomeID tidak ditemukan"})
// 				return
// 			}
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal validasi HomeID"})
// 			return
// 		}
// 	}

// 	// 🔄 Update field yang diperbolehkan
// 	existing.Rating = input.Rating
// 	existing.Note = input.Note
// 	existing.HomeID = input.HomeID

// 	if err := db.DB.Save(&existing).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate testimony"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, existing)
// }
