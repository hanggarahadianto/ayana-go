package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteCashFlow menghapus semua Goods terkait terlebih dahulu, lalu menghapus CashFlow
func DeleteCashFlow(c *gin.Context) {
	cashFlowID := c.Param("id")

	// Cek apakah CashFlow ada
	var cashFlow models.CashFlow
	if err := db.DB.Where("id = ?", cashFlowID).First(&cashFlow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CashFlow not found"})
		return
	}

	// Hapus semua Goods yang terkait dengan CashFlow (jika ada)
	if err := db.DB.Where("cash_flow_id = ?", cashFlowID).Delete(&models.Goods{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated Goods"})
		return
	}

	// Hapus CashFlow setelah Goods dihapus (atau jika tidak ada Goods)
	if err := db.DB.Delete(&cashFlow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete CashFlow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "CashFlow and related Goods deleted successfully"})
}
