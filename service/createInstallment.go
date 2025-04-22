package service

import (
	"ayana/db"
	"ayana/models"
	"errors"
	"time"

	"github.com/google/uuid"
)

func CreateInstallment(input models.JournalEntry, trxCategory models.TransactionCategory) ([]models.JournalEntry, error) {
	var journals []models.JournalEntry

	installmentCount := input.Installment
	if installmentCount <= 0 {
		return nil, errors.New("installment count must be greater than 0")
	}

	installmentAmount := input.Amount / int64(installmentCount)

	for i := 0; i < installmentCount; i++ {
		debitAmount := int64(0)
		if i == 0 {
			debitAmount = input.Amount
		}

		journal := models.JournalEntry{
			ID:                    uuid.New(),
			Invoice:               input.Invoice,
			Description:           input.Description,
			TransactionCategoryID: input.TransactionCategoryID,
			Amount:                installmentAmount,
			Partner:               input.Partner,
			TransactionType:       input.TransactionType,
			Status:                input.Status,
			DateInputed:           input.DateInputed,
			DueDate:               input.DueDate,
			Note:                  input.Note,
			CompanyID:             input.CompanyID,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		journal.Lines = []models.JournalLine{
			{
				ID:          uuid.New(),
				JournalID:   journal.ID,
				AccountID:   trxCategory.DebitAccountID,
				CompanyID:   journal.CompanyID,
				Debit:       debitAmount,
				Credit:      0,
				Description: input.Description,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New(),
				JournalID:   journal.ID,
				AccountID:   trxCategory.CreditAccountID,
				CompanyID:   journal.CompanyID,
				Debit:       0,
				Credit:      installmentAmount,
				Description: input.Description,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := db.DB.Create(&journal).Error; err != nil {
			return nil, err
		}

		journals = append(journals, journal)
	}

	return journals, nil
}
