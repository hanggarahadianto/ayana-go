package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPresenceRules(c *gin.Context) {
	companyIDParam := c.Query("company_id")
	if companyIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'company_id' wajib diisi"})
		return
	}

	companyID, err := uuid.Parse(companyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format 'company_id' tidak valid"})
		return
	}

	var rules []models.PresenceRule
	if err := db.DB.Where("company_id = ?", companyID).Find(&rules).Error; err != nil {
		log.Println("🔴 Gagal mengambil presence rules:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aturan presensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Berhasil mengambil data aturan presensi",
		"presenceRules": rules,
	})
}
