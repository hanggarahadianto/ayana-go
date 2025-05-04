package controller

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetJournalEntriesByCategory(c *gin.Context) {
	// Ambil query params
	companyIDStr := c.DefaultQuery("company_id", "")
	transactionCategoryIDStr := c.Query("transaction_category_id")

	// Validasi Company ID
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	// Validasi Transaction Category ID
	transactionCategoryID, valid := helper.ValidateAndParseTransactionCategoryID(transactionCategoryIDStr, c)
	if !valid {
		return
	}

	// Ambil paginasi
	pagination := helper.GetPagination(c)

	// Query untuk mengambil journal entries dengan relasi "Lines"
	var journalEntries []models.JournalEntry
	err := db.DB.
		Preload("Lines").
		Where("transaction_category_id = ? AND company_id = ?", transactionCategoryID, companyID).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&journalEntries).Error

	// Handle error query
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Journal entries not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch journal entries"})
		return
	}

	// Mapping data ke DTO response
	var responseData []dto.JournalEntryResponse
	for _, entry := range journalEntries {
		var lines []dto.JournalEntryLineItem
		for _, line := range entry.Lines {
			lines = append(lines, dto.JournalEntryLineItem{
				ID:             line.ID.String(),        // Mengonversi UUID ke string
				JournalEntryID: entry.ID.String(),       // Menambahkan ID dari JournalEntry
				AccountID:      line.AccountID.String(), // Jika line.AccountID adalah UUID
				AccountName:    line.Account.Name,       // Mengambil nama akun dari relasi
				Debit:          float64(line.Debit),     // Konversi ke float64
				Credit:         float64(line.Credit),    // Konversi ke float64
				Description:    line.Description,
			})
		}

		// Mengonversi DateInputed jika tidak nil
		var dateInputed time.Time
		if entry.DateInputed != nil {
			dateInputed = *entry.DateInputed // Jika tidak nil, ambil nilai waktu
		}

		// Menambahkan entry ke response data
		responseData = append(responseData, dto.JournalEntryResponse{
			ID:                    entry.ID.String(), // Mengonversi UUID ke string
			Invoice:               entry.Invoice,
			Description:           entry.Description,
			TransactionID:         entry.Transaction_ID,
			TransactionCategoryID: entry.TransactionCategoryID.String(), // Mengonversi UUID ke string
			Amount:                float64(entry.Amount),                // Konversi dari int64 ke float64
			Partner:               entry.Partner,
			TransactionType:       string(entry.TransactionType), // Konversi enum ke string
			Status:                string(entry.Status),          // Konversi enum ke string
			CompanyID:             entry.CompanyID.String(),      // Mengonversi UUID ke string
			DateInputed:           dateInputed,                   // Jika nil, kosongkan
			IsRepaid:              entry.IsRepaid,
			Installment:           entry.Installment,
			Note:                  entry.Note,
			Lines:                 lines,
		})
	}

	// Mengirimkan response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   pagination.Page,
		"limit":  pagination.Limit,
		"data":   responseData,
	})
}
