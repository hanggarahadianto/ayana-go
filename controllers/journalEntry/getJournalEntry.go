package controller

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetJournalEntriesByCategory(c *gin.Context) {
	companyIDStr := c.DefaultQuery("company_id", "")

	// Validasi Company ID
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	pagination := helper.GetPagination(c)

	var total int64
	if err := db.DB.
		Model(&models.JournalEntry{}).
		Where("company_id = ?", companyID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total data"})
		return
	}

	var journalEntries []models.JournalEntry
	err := db.DB.
		Preload("Lines.Account").
		Preload("TransactionCategory").
		Where("company_id = ?", companyID).
		Order("journal_entries.date_inputed DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&journalEntries).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Journal entries tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil journal entries"})
		return
	}

	var responseData []dto.JournalEntryResponse
	for _, entry := range journalEntries {
		responseData = append(responseData, dto.JournalEntryResponse{
			ID:                      entry.ID.String(),
			TransactionID:           entry.Transaction_ID,
			TransactionCategoryID:   entry.TransactionCategoryID.String(),
			TransactionCategoryName: entry.TransactionCategory.Name,
			Invoice:                 entry.Invoice,
			DebitCategory:           entry.TransactionCategory.DebitCategory,
			CreditCategory:          entry.TransactionCategory.CreditCategory,
			Partner:                 entry.Partner,
			Description:             entry.Description,
			Amount:                  entry.Amount,
			TransactionType:         string(entry.TransactionType),
			DebitAccountType:        entry.DebitAccountType,
			CreditAccountType:       entry.CreditAccountType,
			Status:                  string(entry.Status),
			CompanyID:               entry.CompanyID.String(),
			DateInputed:             entry.DateInputed,
			DueDate:                 entry.DueDate,
			RepaymentDate:           entry.RepaymentDate,
			IsRepaid:                entry.IsRepaid,
			Installment:             entry.Installment,
			Note:                    entry.Note,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"journalEntryList": responseData,
			"page":             pagination.Page,
			"limit":            pagination.Limit,
			"total":            total,
		},
		"message": "Data journal entries berhasil diambil",
		"status":  "sukses",
	})
}
