package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// controllers/presence_rule_controller.go

func CreatePresenceRules(c *gin.Context) {
	var input models.PresenceRule

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("🔴 Bind JSON error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format payload salah. Pastikan semua field terisi dengan benar.",
		})
		return
	}

	if input.Day == "" || input.CompanyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'day' dan 'company_id' wajib diisi"})
		return
	}

	if err := db.DB.Create(&input).Error; err != nil {
		log.Println("🔴 Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan data aturan presensi",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Aturan presensi berhasil disimpan",
		"data":    input,
	})
}
