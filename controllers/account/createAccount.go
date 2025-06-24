package controllers

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var input models.Account

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !helper.ValidateCompanyExist(input.CompanyID, c) {
		return
	}

	if err := helper.ValidateAccount(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi apakah kode akun sudah ada
	// Validasi apakah kode akun sudah ada untuk perusahaan yang sama
	var existingAccount models.Account
	if err := db.DB.Where("code = ? AND company_id = ?", input.Code, input.CompanyID).First(&existingAccount).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kode sudah tersedia untuk perusahaan ini",
		})
		return
	}

	account := models.Account{
		Code:        input.Code,
		Name:        input.Name,
		Type:        input.Type,
		Category:    input.Category,
		Description: input.Description,
		CompanyID:   input.CompanyID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   account,
	})
}
