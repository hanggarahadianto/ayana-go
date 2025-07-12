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
	// Ambil ID dari path param
	idParam := c.Param("id")
	ruleID, err := uuid.Parse(idParam)
	if err != nil {
		log.Println("🔴 Invalid UUID:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid.",
		})
		return
	}

	var input models.PresenceRule

	// Bind body JSON ke struct
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

	// Cek apakah rule dengan ID tersebut ada
	var existing models.PresenceRule
	if err := db.DB.Where("id = ?", ruleID).First(&existing).Error; err != nil {
		log.Println("🔴 Rule tidak ditemukan:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Aturan presensi tidak ditemukan.",
		})
		return
	}

	// Update data dari input
	existing.Day = strings.ToLower(input.Day)
	existing.CompanyID = input.CompanyID
	existing.IsHoliday = input.IsHoliday
	existing.StartTime = input.StartTime
	existing.EndTime = input.EndTime
	existing.GracePeriodMins = input.GracePeriodMins
	existing.ArrivalTolerances = input.ArrivalTolerances
	existing.DepartureTolerances = input.DepartureTolerances

	// Simpan ke DB
	if err := db.DB.Save(&existing).Error; err != nil {
		log.Println("🔴 Gagal update:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate aturan presensi.",
		})
		return
	}

	// Sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Aturan presensi berhasil diperbarui.",
		"data":    existing,
	})
}
