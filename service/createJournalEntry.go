package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func ProcessSingleJournalEntry(input models.JournalEntry) (models.JournalEntry, error) {
	return createJournalEntryService(input)
}

func ProcessMultipleJournalEntries(inputs []models.JournalEntry) ([]models.JournalEntry, error) {
	var results []models.JournalEntry

	for _, input := range inputs {
		if input.ID != uuid.Nil {
			// Update journal existing
			err := updateJournalStatus(input.ID)
			if err != nil {
				return nil, err
			}

			var updatedJournal models.JournalEntry
			if err := db.DB.Preload("Lines.Account").
				Preload("TransactionCategory.DebitAccount").
				Preload("TransactionCategory.CreditAccount").
				First(&updatedJournal, "id = ?", input.ID).Error; err != nil {
				return nil, err
			}

			results = append(results, updatedJournal)
		} else {
			// Create new journal
			entry, err := createJournalEntryService(input)
			if err != nil {
				return nil, err
			}
			results = append(results, entry)
		}
	}

	return results, nil
}

func createJournalEntryService(input models.JournalEntry) (models.JournalEntry, error) {
	if input.TransactionCategoryID == uuid.Nil || input.Amount <= 0 || input.CompanyID == uuid.Nil {
		return models.JournalEntry{}, fmt.Errorf("missing required fields")
	}

	tx := db.DB.Begin() // START transaction
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var company models.Company
	if err := tx.First(&company, "id = ?", input.CompanyID).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	var trxCategory models.TransactionCategory
	if err := tx.Preload("DebitAccount").
		Preload("CreditAccount").
		First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	now := time.Now()
	journalID := uuid.New()

	if input.Installment > 0 {
		tx.Rollback() // stop transaction; use another flow for installment
		installmentJournals, err := CreateInstallmentJournals(input)
		if err != nil {
			return models.JournalEntry{}, err
		}
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
		RepaymentDate:         input.RepaymentDate,
		DueDate:               input.DueDate,
		Note:                  input.Note,
		CompanyID:             input.CompanyID,
		CreatedAt:             now,
		UpdatedAt:             now,
		Lines: []models.JournalLine{
			{
				ID:                uuid.New(),
				JournalID:         journalID,
				AccountID:         trxCategory.DebitAccountID,
				CompanyID:         input.CompanyID,
				Debit:             input.Amount,
				Credit:            0,
				DebitAccountType:  trxCategory.DebitAccountType,
				CreditAccountType: trxCategory.CreditAccountType,
				TransactionType:   input.TransactionType,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                uuid.New(),
				JournalID:         journalID,
				AccountID:         trxCategory.CreditAccountID,
				CompanyID:         input.CompanyID,
				Debit:             0,
				Credit:            input.Amount,
				DebitAccountType:  trxCategory.DebitAccountType,
				CreditAccountType: trxCategory.CreditAccountType,
				TransactionType:   input.TransactionType,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}

	if err := tx.Create(&journal).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	var journalWithDetails models.JournalEntry
	if err := tx.Preload("Lines.Account").
		Preload("TransactionCategory.DebitAccount").
		Preload("TransactionCategory.CreditAccount").
		First(&journalWithDetails, "id = ?", journal.ID).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	// Commit terlebih dahulu, baru index (karena indexing bukan bagian dari DB transaction)
	tx.Commit()

	// Index ke Typesense setelah commit (karena indexing bukan atomic DB op)
	if err := IndexJournalDocument(journalWithDetails); err != nil {
		return journalWithDetails, fmt.Errorf("data saved but indexing failed: %w", err)
	}

	return journalWithDetails, nil
}
