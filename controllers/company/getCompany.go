package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCompany(c *gin.Context) {
	// Ambil query params untuk pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var companyList []models.Company
	var total int64

	// Hitung total data
	if err := db.DB.Model(&models.Company{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghitung data",
			"error":   err.Error(),
		})
		return
	}

	// Ambil data dengan limit & offset
	if err := db.DB.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&companyList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data perusahaan",
			"error":   err.Error(),
		})
		return
	}

	// Response dengan format yang kamu minta
	c.JSON(http.StatusOK, gin.H{
		"data":       companyList,
		"page":       page,
		"limit":      limit,
		"total_data": total,
		"total_page": (total + int64(limit) - 1) / int64(limit),
		"status":     "success",
	})
}
