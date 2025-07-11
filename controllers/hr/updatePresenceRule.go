package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdatePresenceRule handles PUT /presence-rules/:id
func UpdatePresenceRule(c *gin.Context) {
	id := c.Param("id")
	ruleID, err := uuid.Parse(id)
	if err != nil {
		log.Println("🔴 Invalid UUID:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid.",
		})
		return
	}

	var input models.PresenceRule
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("🔴 Bind JSON error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format payload salah. Pastikan semua field diisi dengan benar.",
		})
		return
	}

	// Validasi field penting
	if strings.TrimSpace(input.Day) == "" || input.CompanyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Field 'day' dan 'company_id' wajib diisi.",
		})
		return
	}

	// Cek apakah rule-nya ada
	var existing models.PresenceRule
	if err := db.DB.First(&existing, "id = ?", ruleID).Error; err != nil {
		log.Println("🔴 Rule tidak ditemukan:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Aturan presensi tidak ditemukan.",
		})
		return
	}

	// Update data
	existing.Day = strings.ToLower(input.Day)
	existing.CompanyID = input.CompanyID
	existing.IsHoliday = input.IsHoliday
	existing.StartTime = input.StartTime
	existing.EndTime = input.EndTime
	existing.GracePeriodMins = input.GracePeriodMins
	existing.ArrivalTolerances = input.ArrivalTolerances
	existing.DepartureTolerances = input.DepartureTolerances

	if err := db.DB.Save(&existing).Error; err != nil {
		log.Println("🔴 Gagal update:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate aturan presensi.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Aturan presensi berhasil diperbarui.",
		"data":    existing,
	})
}
