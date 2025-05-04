package controllers

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	accountType := c.Query("type")
	category := c.Query("category")
	categoryOnly := c.Query("category_only") == "true" // ✅ Tambahan filter baru

	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	pagination := helper.GetPagination(c)

	if categoryOnly {
		var categories []struct {
			ID       string `json:"id"`
			Category string `json:"category"`
		}

		if err := db.DB.Model(&models.Account{}).
			Select("DISTINCT ON (category) id, category").
			Where("company_id = ?", companyID).
			Order("category").
			Scan(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   categories,
		})
		return
	}

	var accounts []models.Account
	query := db.DB.Model(&models.Account{}).Where("company_id = ?", companyID)

	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}
	if category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", category)
	}

	if err := query.
		Order("code ASC").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	var total int64
	countQuery := db.DB.Model(&models.Account{}).Where("company_id = ?", companyID)
	if accountType != "" {
		countQuery = countQuery.Where("type = ?", accountType)
	}
	if category != "" {
		countQuery = countQuery.Where("category = ?", category)
	}
	countQuery.Count(&total)

	var responseData []dto.AccountResponse
	for _, a := range accounts {
		responseData = append(responseData, dto.AccountResponse{
			ID:          a.ID,
			Code:        a.Code,
			Name:        a.Name,
			Type:        a.Type,
			Category:    a.Category,
			Description: a.Description,
			CompanyID:   a.CompanyID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   pagination.Page,
		"limit":  pagination.Limit,
		"total":  total,
		"data":   responseData,
	})
}
