package controller

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOutstandingDebts(c *gin.Context) {
	companyID := c.Query("company_id")
	status := c.Query("status")
	if companyID == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	if c.Query("summary_only") == "true" {
		totalOutstandingDebt, err := service.GetOutstandingDebtSummaryOnly(companyID, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate summary total"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"total_outstandingDebt": totalOutstandingDebt,
			},
			"message": "OutstandingDebt summary retrieved successfully",
			"status":  "success",
		})

		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var entries []models.JournalEntry
	var total int64
	var totalDebt int64

	// Subquery untuk filter hanya journal yang memiliki credit > 0 (hutang)
	subQuery := db.DB.
		Model(&models.JournalLine{}).
		Select("journal_id").
		Where("credit > 0 AND company_id = ?", companyID)

	// Hitung total data

	if err := db.DB.Model(&models.JournalEntry{}).
		Where("id IN (?) AND status = ? AND is_repaid = false", subQuery, status).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"debtList":   []dto.JournalEntryResponse{},
				"total_debt": 0,
				"page":       page,
				"limit":      limit,
				"total":      0,
			},
			"message": "No data available",
			"status":  "success",
		})
		return
	}

	// Hitung total utang (SUM amount)
	if err := db.DB.Model(&models.JournalEntry{}).
		Where("id IN (?) AND status = ? AND is_repaid = false", subQuery, status).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalDebt).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"debtList":   []dto.JournalEntryResponse{},
				"total_debt": 0,
				"page":       page,
				"limit":      limit,
				"total":      total,
			},
			"message": "No data available",
			"status":  "success",
		})
		return
	}

	// Ambil data dengan pagination
	if err := db.DB.
		Where("id IN (?) AND status = ? AND is_repaid = false", subQuery, status).
		Order("due_date ASC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch data"})
		return
	}

	// Buat response
	var responseData []dto.JournalEntryResponse
	for _, entry := range entries {
		responseData = append(responseData, dto.JournalEntryResponse{
			ID:                    entry.ID.String(),
			Invoice:               entry.Invoice,
			Description:           entry.Description,
			TransactionCategoryID: entry.TransactionCategoryID.String(),
			Amount:                float64(entry.Amount),
			Partner:               entry.Partner,
			TransactionType:       string(entry.TransactionType),
			Status:                string(entry.Status),
			CompanyID:             entry.CompanyID.String(),
			DateInputed:           *entry.DateInputed,
			DueDate:               *entry.DueDate,
			IsRepaid:              entry.IsRepaid,
			Installment:           entry.Installment,
			Note:                  entry.Note,
			Lines:                 nil,
		})
	}

	// Kirim JSON response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"debtList":              responseData,
			"total_outstandingDebt": totalDebt,
			"page":                  page,
			"limit":                 limit,
			"total":                 total,
		},
		"message": "Outstanding debts retrieved successfully",
		"status":  "success",
	})
}
