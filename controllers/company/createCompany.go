package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateCompany(c *gin.Context) {
	var companyData models.Company
	username, _ := c.Get("username")
	if username != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Hanya superadmin yang dapat mengakses data ini",
			"status":  "error",
		})
		return
	}

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&companyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Cek apakah CompanyCode sudah ada
	var existingCompany models.Company
	if err := db.DB.Where("company_code = ?", companyData.CompanyCode).First(&existingCompany).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Code already exist",
		})
		return
	}

	// Buat record baru
	now := time.Now()
	newCompany := models.Company{
		Title:       companyData.Title,
		CompanyCode: companyData.CompanyCode,
		HasProduct:  companyData.HasProduct,
		HasProject:  companyData.HasProject,
		HasCustomer: companyData.HasCustomer,
		IsRetail:    companyData.IsRetail,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.DB.Create(&newCompany).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newCompany,
	})
}
