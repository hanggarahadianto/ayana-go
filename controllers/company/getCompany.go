package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCompany(c *gin.Context) {
	var companyList []models.Company

	// Query database dengan pengecekan lebih baik
	err := db.DB.
		Order("created_at DESC, updated_at DESC").
		Find(&companyList).Error

	if err != nil {
		// Jika terjadi error database
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to query database",
			"error":   err.Error(),
		})
		return
	}

	// Jika tidak ada data yang ditemukan
	if len(companyList) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No companies found",
			"data":    []models.Company{},
		})
		return
	}

	// Jika data ditemukan
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Companies retrieved successfully",
		"data":    companyList,
	})
}
