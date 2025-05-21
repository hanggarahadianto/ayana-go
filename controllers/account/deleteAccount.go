package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeleteAccount menghapus akun dan relasi TransactionCategory terkait
func DeleteAccount(c *gin.Context) {
	accountIDParam := c.Param("id")
	accountID, err := uuid.Parse(accountIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Periksa apakah akun ada
	var account models.Account
	if err := db.DB.First(&account, "id = ?", accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Hapus semua TransactionCategory yang memiliki akun ini sebagai debit atau credit
	if err := db.DB.Where("debit_account_id = ? OR credit_account_id = ?", accountID, accountID).
		Delete(&models.TransactionCategory{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related transaction categories"})
		return
	}

	// Hapus akun itu sendiri
	if err := db.DB.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account and related transaction categories deleted successfully",
	})
}
