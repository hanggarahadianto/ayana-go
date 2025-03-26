package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCashFlow(c *gin.Context) {
	cashFlowID := c.Param("id")
	var cashFlow models.CashFlow

	// Cek apakah CashFlow ada
	if err := db.DB.Where("id = ?", cashFlowID).First(&cashFlow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CashFlow not found"})
		return
	}

	// Cek apakah ada Goods yang terkait
	var goodsCount int64
	db.DB.Model(&models.Goods{}).Where("cash_flow_id = ?", cashFlowID).Count(&goodsCount)
	if goodsCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete CashFlow with associated Goods"})
		return
	}

	// Hapus CashFlow jika tidak ada Goods terkait
	if err := db.DB.Delete(&cashFlow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete CashFlow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
