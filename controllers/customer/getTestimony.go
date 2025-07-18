package controllers

import (
	"ayana/db"
	lib "ayana/lib"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTestimonies(c *gin.Context) {
	pagination := lib.GetPagination(c)
	if !lib.ValidatePagination(pagination, c) {
		return
	}

	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id wajib diisi"})
		return
	}

	var (
		testimonies []models.Testimony
		total       int64
	)

	// Hitung total testimony untuk perusahaan
	if err := db.DB.Model(&models.Testimony{}).
		Where("company_id = ?", companyID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghitung data testimony: " + err.Error(),
		})
		return
	}

	// Ambil data testimony dengan relasi Customer dan Home
	if err := db.DB.
		Preload("Customer").
		Preload("Customer.Home").
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset((pagination.Page - 1) * pagination.Limit).
		Find(&testimonies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data testimony: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Data testimony berhasil diambil",
		"data": gin.H{
			"testimonyList":   testimonies,
			"total_testimony": total,
			"page":            pagination.Page,
			"limit":           pagination.Limit,
			"total":           total,
		},
	})
}
