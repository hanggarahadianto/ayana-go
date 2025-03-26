package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateGoods membuat beberapa barang baru dalam database
func CreateGood(c *gin.Context) {
	var goods []models.Goods

	// Bind JSON input ke slice of Goods
	if err := c.ShouldBindJSON(&goods); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mulai transaksi database
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		}
	}()

	// Iterasi melalui setiap barang dalam array
	for i := range goods {
		// Validasi apakah cash_flow_id ada di tabel CashFlow
		var cashFlow models.CashFlow
		if err := tx.First(&cashFlow, "id = ?", goods[i].CashFlowId).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Cashflow not found for good"})
			return
		}

		// Set ID jika belum ada
		if goods[i].ID == uuid.Nil {
			goods[i].ID = uuid.New()
		}

		// Simpan barang ke database
		if err := tx.Create(&goods[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Good "})
			return
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   goods,
	})
}
