package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) {
	companyID := c.Query("company_id")
	accountType := c.Query("type")

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

	var accounts []models.Account
	query := db.DB.Model(&models.Account{}).Where("company_id = ?", companyID)

	// Apply filter if type is provided
	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}

	// Fetch paginated data
	if err := query.
		Order("code ASC").
		Limit(limit).
		Offset(offset).
		Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	// Count total records without limit and offset
	var total int64
	countQuery := db.DB.Model(&models.Account{}).Where("company_id = ?", companyID)
	if accountType != "" {
		countQuery = countQuery.Where("type = ?", accountType)
	}
	countQuery.Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   accounts,
	})
}
