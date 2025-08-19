package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateCompany(c *gin.Context) {
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

	// Pastikan company yang mau diupdate ada
	var existingCompany models.Company
	if err := db.DB.First(&existingCompany, companyData.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Company not found",
		})
		return
	}

	// Cek apakah company_code sudah dipakai oleh company lain
	var companyWithSameCode models.Company
	if err := db.DB.Where("company_code = ? AND id <> ?", companyData.CompanyCode, companyData.ID).
		First(&companyWithSameCode).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Code already exist",
		})
		return
	}

	// Update data
	existingCompany.Title = companyData.Title
	existingCompany.HasProduct = companyData.HasProduct
	existingCompany.HasCustomer = companyData.HasCustomer
	existingCompany.HasProject = companyData.HasProject
	existingCompany.IsRetail = companyData.IsRetail
	existingCompany.CompanyCode = companyData.CompanyCode
	existingCompany.UpdatedAt = time.Now()

	if err := db.DB.Save(&existingCompany).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existingCompany,
	})
}
