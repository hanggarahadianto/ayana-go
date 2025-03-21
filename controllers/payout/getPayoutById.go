package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPayoutsByCompany(c *gin.Context) {
	// Ambil query parameter dari request
	companyID := c.Query("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default page = 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit = 10

	// Validasi jika company_id kosong
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	// Konversi paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query untuk mengambil data payout berdasarkan company_id dengan paginasi
	var payouts []models.Payout
	if err := db.DB.Where("company_id = ?", companyID).
		Limit(limit).
		Offset(offset).
		Find(&payouts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payouts"})
		return
	}

	// Hitung total data payout untuk company_id tersebut
	var total int64
	db.DB.Model(&models.Payout{}).Where("company_id = ?", companyID).Count(&total)

	// Kirim response dengan paginasi
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   payouts,
	})
}
