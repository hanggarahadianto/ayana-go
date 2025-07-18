package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DeleteTestimony(c *gin.Context) {
	id := c.Param("id")

	// ✅ Parse UUID
	testimonyID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// 🔍 Cek apakah testimony ada
	var testimony models.Testimony
	if err := db.DB.First(&testimony, "id = ?", testimonyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Testimony tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data testimony"})
		return
	}

	// 🗑️ Hapus testimony
	if err := db.DB.Delete(&testimony).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus testimony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Testimony berhasil dihapus"})
}
