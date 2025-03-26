package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UpdateGood menangani pembuatan, pembaruan, dan penghapusan data Goods berdasarkan payload
func UpdateGood(c *gin.Context) {
	var goodsPayload []models.Goods

	// Bind JSON payload ke struct goodsPayload
	if err := c.ShouldBindJSON(&goodsPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response []models.Goods
	var errors []error

	// Loop melalui setiap item dalam payload
	for _, good := range goodsPayload {
		goodID, err := uuid.Parse(good.ID.String())

		// Jika ID tidak valid atau kosong, buat record baru
		if err != nil || good.ID == uuid.Nil {
			newGood := models.Goods{
				ID:         good.ID,
				GoodsName:  good.GoodsName,
				Status:     good.Status,
				Quantity:   good.Quantity,
				CostsDue:   good.CostsDue,
				Price:      good.Price,
				Unit:       good.Unit,
				TotalCost:  good.TotalCost,
				CashFlowId: good.CashFlowId,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := db.DB.Create(&newGood).Error; err != nil {
				errors = append(errors, err)
				continue
			}
			response = append(response, newGood)
		} else {
			// Jika ID valid, lakukan update
			var existingGood models.Goods
			result := db.DB.First(&existingGood, "id = ?", goodID)

			// Jika data tidak ditemukan, tambahkan ke errors
			if result.Error == gorm.ErrRecordNotFound {
				errors = append(errors, result.Error)
				continue
			} else if result.Error != nil {
				errors = append(errors, result.Error)
				continue
			}

			// Update field dengan nilai baru
			existingGood.GoodsName = good.GoodsName
			existingGood.Status = good.Status
			existingGood.Quantity = good.Quantity
			existingGood.CostsDue = good.CostsDue
			existingGood.Price = good.Price
			existingGood.Unit = good.Unit
			existingGood.TotalCost = good.TotalCost
			existingGood.CashFlowId = good.CashFlowId
			existingGood.UpdatedAt = time.Now()

			if err := db.DB.Save(&existingGood).Error; err != nil {
				errors = append(errors, err)
				continue
			}
			response = append(response, existingGood)
		}
	}

	// Jika ada error selama proses, kirimkan HTTP 206 Partial Content
	if len(errors) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message": "Some operations failed",
			"data":    response,
			"errors":  errors,
		})
		return
	}

	// Jika semua operasi berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "Operations completed successfully",
		"data":    response,
	})
}
