package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCompany(c *gin.Context) {
	// Ambil query parameter dari request

	// Query untuk mencari company berdasarkan company_id
	var companyList []models.Company

	if err := db.DB.Debug().Order("created_at desc, updated_at desc").Find(&companyList).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// Kirim response jika ditemukan
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   companyList,
	})
}
