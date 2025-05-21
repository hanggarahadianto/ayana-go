package controllers

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EditAccount mengedit data akun berdasarkan ID
func EditAccount(c *gin.Context) {
	accountIDParam := c.Param("id")
	accountID, err := uuid.Parse(accountIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var input models.Account
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi company
	if !helper.ValidateCompanyExist(input.CompanyID, c) {
		return
	}

	// Validasi data account
	if err := helper.ValidateAccount(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari account yang ingin diupdate
	var existingAccount models.Account
	if err := db.DB.First(&existingAccount, "id = ?", accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Cek jika code berubah, pastikan tidak duplikat
	if input.Code != existingAccount.Code {
		var otherAccount models.Account
		if err := db.DB.Where("code = ?", input.Code).First(&otherAccount).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Kode sudah tersedia"})
			return
		}
	}

	// Update data
	existingAccount.Code = input.Code
	existingAccount.Name = input.Name
	existingAccount.Type = input.Type
	existingAccount.Category = input.Category
	existingAccount.Description = input.Description
	existingAccount.CompanyID = input.CompanyID
	existingAccount.UpdatedAt = time.Now()

	if err := db.DB.Save(&existingAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existingAccount,
	})
}
