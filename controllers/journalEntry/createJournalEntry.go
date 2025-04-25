package controller

import (
	"ayana/db"
	"ayana/models"
	"ayana/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateJournalEntry(c *gin.Context) {
	var input models.JournalEntry

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "details": err.Error()})
		return
	}

	// 2. Validasi minimal input
	if input.TransactionCategoryID == uuid.Nil || input.Amount <= 0 || input.CompanyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing required fields"})
		return
	}

	// 3. Validasi company
	var company models.Company
	if err := db.DB.First(&company, "id = ?", input.CompanyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Company does not exist"})
		return
	}

	// 4. Ambil kategori transaksi lengkap dengan akun
	var trxCategory models.TransactionCategory
	if err := db.DB.
		Preload("DebitAccount").
		Preload("CreditAccount").
		First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Transaction category not found"})
		return
	}

	// 5. Handle installment jika ada
	if input.Installment > 0 {
		journals, err := service.CreateInstallmentJournals(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create installment journals", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": journals})
		return
	}

	// 6. Buat journal entry biasa
	now := time.Now()
	journal := models.JournalEntry{
		ID:                    uuid.New(),
		Invoice:               input.Invoice,
		Description:           input.Note,
		TransactionCategoryID: input.TransactionCategoryID,
		Amount:                input.Amount,
		DebitAccountType:      trxCategory.DebitAccount.Type,
		CreditAccountType:     trxCategory.CreditAccount.Type,
		Partner:               input.Partner,
		TransactionType:       input.TransactionType,
		Status:                input.Status,
		IsRepaid:              input.IsRepaid,
		DateInputed:           input.DateInputed,
		DueDate:               input.DueDate,
		Note:                  input.Note,
		CompanyID:             input.CompanyID,
		CreatedAt:             now,
		UpdatedAt:             now,
		Lines: []models.JournalLine{
			{
				ID:          uuid.New(),
				JournalID:   input.ID,
				AccountID:   trxCategory.DebitAccountID,
				CompanyID:   input.CompanyID,
				Debit:       input.Amount,
				Credit:      0,
				Description: input.Description,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          uuid.New(),
				JournalID:   input.ID,
				AccountID:   trxCategory.CreditAccountID,
				CompanyID:   input.CompanyID,
				Debit:       0,
				Credit:      input.Amount,
				Description: input.Description,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}

	// 7. Simpan journal + lines
	if err := db.DB.Create(&journal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create journal entry", "details": err.Error()})
		return
	}

	// 8. Ambil data lengkap setelah insert
	var journalWithDetails models.JournalEntry
	if err := db.DB.Preload("Lines.Account").
		Preload("TransactionCategory.DebitAccount").
		Preload("TransactionCategory.CreditAccount").
		First(&journalWithDetails, "id = ?", journal.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to preload journal details", "details": err.Error()})
		return
	}

	// 9. Return success
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": journalWithDetails})
}
