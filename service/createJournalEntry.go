package service

import (
	"fmt"
	"time"

	"ayana/db"
	"ayana/models"

	"github.com/google/uuid"
)

func ProcessSingleJournalEntry(input models.JournalEntry) (models.JournalEntry, error) {
	return createJournalEntryService(input)
}

func ProcessMultipleJournalEntries(inputs []models.JournalEntry) ([]models.JournalEntry, error) {
	var results []models.JournalEntry

	for _, input := range inputs {
		entry, err := createJournalEntryService(input)
		if err != nil {
			return nil, err
		}
		results = append(results, entry)
	}

	return results, nil
}

// Core logic create journal
func createJournalEntryService(input models.JournalEntry) (models.JournalEntry, error) {
	if input.TransactionCategoryID == uuid.Nil || input.Amount <= 0 || input.CompanyID == uuid.Nil {
		return models.JournalEntry{}, fmt.Errorf("missing required fields")
	}

	var company models.Company
	if err := db.DB.First(&company, "id = ?", input.CompanyID).Error; err != nil {
		return models.JournalEntry{}, err
	}

	var trxCategory models.TransactionCategory
	if err := db.DB.
		Preload("DebitAccount").
		Preload("CreditAccount").
		First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
		return models.JournalEntry{}, err
	}

	now := time.Now()
	journalID := uuid.New()

	if input.Installment > 0 {
		installmentJournals, err := CreateInstallmentJournals(input)
		if err != nil {
			return models.JournalEntry{}, err
		}

		// Kembalikan journal pertama yang dibuat sebagai hasil
		return installmentJournals[0], nil
	}

	journal := models.JournalEntry{
		ID:                    journalID,
		Transaction_ID:        input.Transaction_ID,
		Invoice:               input.Invoice,
		Description:           input.Note,
		TransactionCategoryID: input.TransactionCategoryID,
		Amount:                input.Amount,
		DebitAccountType:      trxCategory.DebitAccountType,
		CreditAccountType:     trxCategory.CreditAccountType,
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
				JournalID:   journalID,
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
				JournalID:   journalID,
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

	if err := db.DB.Create(&journal).Error; err != nil {
		return models.JournalEntry{}, err
	}

	var journalWithDetails models.JournalEntry
	if err := db.DB.Preload("Lines.Account").
		Preload("TransactionCategory.DebitAccount").
		Preload("TransactionCategory.CreditAccount").
		First(&journalWithDetails, "id = ?", journal.ID).Error; err != nil {
		return models.JournalEntry{}, err
	}

	return journalWithDetails, nil
}
