package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateGood(c *gin.Context) {
	var goodsPayload []models.Goods

	// Bind JSON payload ke struct
	if err := c.ShouldBindJSON(&goodsPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response []models.Goods
	var errors []error

	// Ambil semua ID yang terkait dengan cash_flow_id yang dikirimkan dalam payload
	var cashFlowIDs []uuid.UUID
	for _, good := range goodsPayload {
		cashFlowIDs = append(cashFlowIDs, good.CashFlowId)
	}

	var existingGoods []models.Goods
	if err := db.DB.Where("cash_flow_id IN (?)", cashFlowIDs).Find(&existingGoods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing goods"})
		return
	}

	existingGoodsMap := make(map[uuid.UUID]models.Goods)
	for _, g := range existingGoods {
		existingGoodsMap[g.ID] = g
	}

	processedIDs := make(map[uuid.UUID]bool)

	for _, good := range goodsPayload {
		goodID := good.ID
		if goodID == uuid.Nil {
			// Create new record jika ID kosong
			good.ID = uuid.New()
			good.CreatedAt = time.Now()
			good.UpdatedAt = time.Now()

			if err := db.DB.Create(&good).Error; err != nil {
				errors = append(errors, err)
				continue
			}
			response = append(response, good)
		} else {
			// Update existing record jika ID ditemukan
			if existingGood, found := existingGoodsMap[goodID]; found {
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
				processedIDs[goodID] = true
			}
		}
	}

	// Hapus data yang tidak termasuk dalam payload
	for id, existingGood := range existingGoodsMap {
		if !processedIDs[id] {
			if err := db.DB.Delete(&existingGood).Error; err != nil {
				errors = append(errors, err)
			}
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
