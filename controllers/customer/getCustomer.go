package controllers

import (
	"ayana/db"
	"ayana/utils/helper"

	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	pagination := helper.GetPagination(c)

	if !helper.ValidatePagination(pagination, c) {
		return
	}

	var customers []models.Customer
	var total int64

	// Hitung total data
	if err := db.DB.Model(&models.Customer{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count customers"})
		return
	}

	// Ambil data dengan pagination
	if err := db.DB.Preload("Home").
		Order("updated_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  customers,
		"page":  pagination.Page,
		"limit": pagination.Limit,
		"total": total,
	})
}
