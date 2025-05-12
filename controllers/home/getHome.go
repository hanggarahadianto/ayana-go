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

	var homeList []models.Home
	var total int64

	// Hitung total data
	db.DB.Model(&models.Home{}).Count(&total)

	// Ambil data dengan limit dan offset
	result := db.DB.Order("sequence asc").
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
