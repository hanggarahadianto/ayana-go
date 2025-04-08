package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDebtByID(c *gin.Context) {
	companyID := c.Query("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var debts []models.Payout
	if err := db.DB.
		Where("company_id = ? AND status = ?", companyID, "tempo").
		Limit(limit).
		Offset(offset).
		Find(&debts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch debts"})
		return
	}

	var total int64
	db.DB.Model(&models.Payout{}).
		Where("company_id = ? AND status = ?", companyID, "tempo").
		Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   debts,
	})
}
