package controllers

import (
	"ayana/db"
	"ayana/utils/helper"

	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomes(c *gin.Context) {
	pagination := helper.GetPagination(c)

	// Validasi parameter pagination
	if !helper.ValidatePagination(pagination, c) {
		return
	}

	// Ambil filter status dari query
	status := c.Query("status")

	var homeList []models.Home
	var total int64

	// Mulai query DB
	query := db.DB.Model(&models.Home{})

	// Filter berdasarkan status jika ada
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Hitung total data dengan filter
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Count error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghitung data",
		})
		return
	}

	// Ambil data dengan limit, offset, dan order
	result := query.Order("sequence asc").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&homeList)

	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       homeList,
		"page":       pagination.Page,
		"limit":      pagination.Limit,
		"total_data": total,
		"total_page": (total + int64(pagination.Limit) - 1) / int64(pagination.Limit), // pembulatan ke atas
		"status":     "success",
	})
}
