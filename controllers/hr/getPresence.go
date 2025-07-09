package controllers

import (
	"ayana/db"
	lib "ayana/lib"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPresence(c *gin.Context) {
	var presences []models.Presence
	var total int64
	pagination := lib.GetPagination(c)

	if !lib.ValidatePagination(pagination, c) {
		return
	}

	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company Id wajib diisi"})
		return
	}

	search := c.Query("search")

	// 🔁 Build query awal
	query := db.DB.
		Model(&models.Presence{}).
		Preload("Employee"). // preload relasi employee
		Where("company_id = ?", companyID)

	// 🔍 Filter by employee name (dari relasi)
	if search != "" {
		query = query.Joins("JOIN employees ON employees.id = presences.employee_id").
			Where("employees.name ILIKE ?", "%"+search+"%")
	}

	// 🔢 Hitung total
	if err := query.Count(&total).Error; err != nil {
		log.Println("🔴 Count error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data presensi"})
		return
	}

	// 📦 Ambil data
	if err := query.
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Order("scan_date DESC, scan_time DESC").
		Find(&presences).Error; err != nil {
		log.Println("🔴 Find error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data presensi"})
		return
	}

	// ✅ Response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"presenceList":   presences,
			"total_presence": total,
			"page":           pagination.Page,
			"limit":          pagination.Limit,
			"total":          total,
		},
		"message": "Data presensi berhasil diambil",
		"status":  "sukses",
	})
}
