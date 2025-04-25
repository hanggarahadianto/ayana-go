package controller

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func GetJournalEntriesByCategory(c *gin.Context) {
	transactionCategoryID := c.DefaultQuery("transaction_category_id", "")
	companyID := c.DefaultQuery("company_id", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if transactionCategoryID == "" || companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction category ID and Company ID are required"})
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var journalEntries []models.JournalEntry

	err := db.DB.
		Preload("Lines").
		Where("transaction_category_id = ? AND company_id = ?", transactionCategoryID, companyID).
		Limit(limit).
		Offset(offset).
		Find(&journalEntries).Error

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

	// Mapping ke response DTO
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

		// Untuk field `Amount` (int64 ke float64), `TransactionType` (models.TransactionType ke string), `Status` (models.Status ke string), dan `DateInputed` (pointer *time.Time ke time.Time)
		var dateInputed time.Time
		if entry.DateInputed != nil {
			dateInputed = *entry.DateInputed // Jika tidak nil, ambil nilai waktu
		}

		responseData = append(responseData, dto.JournalEntryResponse{
			ID:                    entry.ID.String(), // Mengonversi UUID ke string
			Invoice:               entry.Invoice,
			Description:           entry.Description,
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"data":   responseData,
	})
}
