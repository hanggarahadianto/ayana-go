package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateInstallmentJournals(input models.JournalEntry) ([]models.JournalEntry, error) {
	var trxCategory models.TransactionCategory
	if err := db.DB.Preload("DebitAccount").Preload("CreditAccount").
		First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
		return nil, err
	}

	amountPerInstallment := input.Amount / int64(input.Installment)

	var journals []models.JournalEntry

	for i := 0; i < input.Installment; i++ {
		journalID := uuid.New()
		dueDate := input.DateInputed.AddDate(0, i+1, 0)
		dueDatePtr := &dueDate

		journal := models.JournalEntry{
			ID:                    journalID,
			Invoice:               fmt.Sprintf("%s-%02d", input.Invoice, i+1),
			Description:           fmt.Sprintf("%s - Cicilan %d", input.Description, i+1),
			TransactionCategoryID: input.TransactionCategoryID,
			Amount:                amountPerInstallment,
			Partner:               input.Partner,
			TransactionType:       input.TransactionType,
			Status:                input.Status,
			DateInputed:           input.DateInputed,
			DueDate:               dueDatePtr,
			CompanyID:             input.CompanyID,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		journal.Lines = []models.JournalLine{
			{
				ID:          uuid.New(),
				JournalID:   journalID,
				AccountID:   trxCategory.DebitAccountID,
				CompanyID:   input.CompanyID,
				Debit:       amountPerInstallment,
				Credit:      0,
				Description: "Auto debit cicilan",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New(),
				JournalID:   journalID,
				AccountID:   trxCategory.CreditAccountID,
				CompanyID:   input.CompanyID,
				Debit:       0,
				Credit:      amountPerInstallment,
				Description: "Auto credit cicilan",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		journals = append(journals, journal)
	}

	if err := db.DB.Create(&journals).Error; err != nil {
		return nil, err
	}

	return journals, nil
}
