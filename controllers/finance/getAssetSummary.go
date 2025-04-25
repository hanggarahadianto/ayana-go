package controller

import (
	"ayana/db"
	"ayana/models"
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAssetSummary(c *gin.Context) {
	// Ambil parameter company_id
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	// Mengecek apakah hanya summary yang diminta
	if c.Query("summary_only") == "true" {
		totalAsset, err := service.GetAssetSummaryOnly(companyID.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate summary total"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"total_asset": totalAsset,
			},
			"message": "Asset summary retrieved successfully",
			"status":  "success",
		})
		return
	}

	// Mendapatkan pagination
	pagination := helper.GetPagination(c)

	var total int64
	// Menghitung total asset (penjumlahan amount seluruh data tanpa paginasi)
	if err := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("JOIN transaction_categories ON transaction_categories.id = journal_entries.transaction_category_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ? AND transaction_categories.debit_account_type = ?", companyID, "paid", true, "Asset").
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count assets"})
		return
	}

	var journalLines []models.JournalLine
	// Mengambil data asset dengan paginasi
	if err := db.DB.
		Preload("Account").
		Preload("Journal").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("JOIN transaction_categories ON transaction_categories.id = journal_entries.transaction_category_id").
		Where("journal_entries.company_id = ? AND journal_entries.status = ? AND journal_entries.is_repaid = ? AND transaction_categories.debit_account_type = ?", companyID, "paid", true, "Asset").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&journalLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset data"})
		return
	}

	totalAsset := int64(0)
	assetList := make([]map[string]interface{}, 0)

	// Hindari duplikasi dengan cek ID
	seen := make(map[uuid.UUID]bool)

	for _, line := range journalLines {
		if line.Debit > 0 && !seen[line.ID] {
			seen[line.ID] = true
			totalAsset += line.Debit

			assetList = append(assetList, map[string]interface{}{
				"id":           line.ID,
				"account_name": line.Account.Name,
				"account_type": line.Account.Type,

				"category":    line.Account.Category,
				"description": line.Journal.Note,
				"amount":      line.Debit,
				"date":        line.CreatedAt.Format("2006-01-02"),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"assetList":   assetList,
			"total_asset": totalAsset,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total":       total,
		},
		"message": "Asset summary retrieved successfully",
		"status":  "success",
	})
}
