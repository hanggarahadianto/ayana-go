package hr

import (
	"ayana/db"
	lib "ayana/lib"
	"ayana/models"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

func GetPresenceService(c *gin.Context) (map[string]interface{}, error) {
	var presences []models.Presence
	var total int64
	pagination := lib.GetPagination(c)

	if !lib.ValidatePagination(pagination, c) {
		return nil, errors.New("pagination tidak valid")
	}

	companyID := c.Query("company_id")
	if companyID == "" {
		return nil, errors.New("company ID wajib diisi")
	}

	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	arrivalOnly := c.Query("arrival_only") == "true"
	departureOnly := c.Query("departure_only") == "true"

	// 🔁 Build query awal
	query := db.DB.
		Model(&models.Presence{}).
		Preload("Employee").
		Where("company_id = ?", companyID)

	// 🔍 Filter nama pegawai
	if search != "" {
		query = query.Joins("JOIN employees ON employees.id = presences.employee_id").
			Where("employees.name ILIKE ?", "%"+search+"%")
	}

	// 📅 Filter tanggal
	if startDate != "" && endDate != "" {
		query = query.Where("scan_date BETWEEN ? AND ?", startDate, endDate)
	}

	// 🕒 Filter waktu berangkat (jam < 12:00)
	if arrivalOnly {
		query = query.Where("scan_time < ?", "12:00")
	}

	// 🕒 Filter waktu pulang (jam >= 12:00)
	if departureOnly {
		query = query.Where("scan_time >= ?", "12:00")
	}

	// 🔢 Hitung total data
	if err := query.Count(&total).Error; err != nil {
		log.Println("🔴 Count error:", err)
		return nil, err
	}

	// 📦 Ambil data presensi
	if err := query.
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Order("scan_date DESC, scan_time DESC").
		Find(&presences).Error; err != nil {
		log.Println("🔴 Find error:", err)
		return nil, err
	}

	// ✅ Response
	response := map[string]interface{}{
		"presenceList":   presences,
		"total_presence": total,
		"page":           pagination.Page,
		"limit":          pagination.Limit,
		"total":          total,
	}

	return response, nil
}
