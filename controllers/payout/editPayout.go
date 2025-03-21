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

func EditPayout(c *gin.Context) {
	var cashFlow models.CashFlow

	if err := c.ShouldBindJSON(&cashFlow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCashFlow models.CashFlow
	if err := db.DB.Where("id = ?", cashFlow.ID).Preload("Good").First(&existingCashFlow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cash Flow ID does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Cash Flow"})
		}
		return
	}

	tx := db.DB.Begin()

	existingCashFlow.WeekNumber = cashFlow.WeekNumber
	existingCashFlow.CashIn = cashFlow.CashIn
	existingCashFlow.CashOut = cashFlow.CashOut
	existingCashFlow.Outstanding = cashFlow.Outstanding
	existingCashFlow.ProjectID = cashFlow.ProjectID
	existingCashFlow.UpdatedAt = time.Now()

	if err := tx.Save(&existingCashFlow).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Cash Flow"})
		return
	}

	var goodIDsToDelete []uuid.UUID
	for _, existingGood := range existingCashFlow.Good {
		found := false
		for _, updatedGood := range cashFlow.Good {
			if existingGood.ID == updatedGood.ID {
				found = true
				break
			}
		}
		if !found {
			goodIDsToDelete = append(goodIDsToDelete, existingGood.ID)
		}
	}

	if len(goodIDsToDelete) > 0 {
		if err := tx.Where("id IN ?", goodIDsToDelete).Delete(&models.Goods{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Goods"})
			return
		}
	}

	for _, good := range cashFlow.Good {
		if good.ID != uuid.Nil {
			if err := tx.Model(&models.Goods{}).Where("id = ?", good.ID).Updates(good).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Good"})
				return
			}
		} else {
			good.CashFlowId = existingCashFlow.ID
			if err := tx.Create(&good).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Good"})
				return
			}
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"cash_flow": existingCashFlow,
		},
	})
}
