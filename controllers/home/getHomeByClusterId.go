package controllers

import (
	"ayana/db"
	"ayana/utils/helper"

	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeListByClusterId(c *gin.Context) {
	clusterId := c.Param("cluster_id")
	pagination := helper.GetPagination(c)

	var homes []models.Home
	var total int64

	// Hitung total data
	db.DB.Model(&models.Home{}).Where("cluster_id = ?", clusterId).Count(&total)

	// Ambil data dengan limit dan offset
	result := db.DB.
		Where("cluster_id = ?", clusterId).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order("created_at DESC").
		Find(&homes)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Failed to get homes by cluster ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   homes,
		"total":  total,
		"limit":  pagination.Limit,
		"offset": pagination.Offset,
	})
}
