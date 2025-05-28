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

	var journalEntries []models.JournalEntry
	amountPerInstallment := input.Amount / int64(input.Installment)
	var totalCredit int64

	// Buat jurnal utama untuk debit penuh
	mainJournalID := uuid.New()
	mainJournal := models.JournalEntry{
		ID:                    mainJournalID,
		Invoice:               input.Invoice,
		Transaction_ID:        input.Transaction_ID,
		Description:           input.Description,
		TransactionCategoryID: input.TransactionCategoryID,
		Amount:                input.Amount,
		Partner:               input.Partner,
		TransactionType:       input.TransactionType,
		DebitAccountType:      trxCategory.DebitAccount.Type,
		CreditAccountType:     trxCategory.CreditAccount.Type,
		Status:                "unpaid",
		IsRepaid:              false,
		DateInputed:           input.DateInputed,
		DueDate:               nil,
		Note:                  input.Note,
		CompanyID:             input.CompanyID,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		Lines: []models.JournalLine{
			{
				ID:                uuid.New(),
				JournalID:         mainJournalID,
				AccountID:         trxCategory.DebitAccountID,
				CompanyID:         input.CompanyID,
				Debit:             input.Amount,
				Credit:            0,
				DebitAccountType:  trxCategory.DebitAccountType,
				CreditAccountType: trxCategory.CreditAccountType,
				TransactionType:   input.TransactionType,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	journalEntries = append(journalEntries, mainJournal)

	// Buat journal entry terpisah untuk setiap cicilan kredit
	for i := 0; i < input.Installment; i++ {
		journalID := uuid.New()
		dueDate := input.DateInputed.AddDate(0, i+1, 0)
		dueDatePtr := &dueDate

		amount := amountPerInstallment
		if i == input.Installment-1 {
			amount = input.Amount - totalCredit
		}
		totalCredit += amount

		journal := models.JournalEntry{
			ID:                    journalID,
			Invoice:               fmt.Sprintf("%s-%02d", "Tempo "+input.Invoice, i+1),
			Transaction_ID:        input.Transaction_ID,
			Description:           fmt.Sprintf("%s - Cicilan %d", input.Note, i+1),
			TransactionCategoryID: input.TransactionCategoryID,
			Amount:                amount,
			Partner:               input.Partner,
			TransactionType:       "payout",
			DebitAccountType:      trxCategory.DebitAccount.Type,
			CreditAccountType:     trxCategory.CreditAccount.Type,
			Status:                "unpaid",
			IsRepaid:              false,
			DateInputed:           input.DateInputed,
			DueDate:               dueDatePtr,
			Note:                  input.Note,
			CompanyID:             input.CompanyID,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			Lines: []models.JournalLine{
				{
					ID:                uuid.New(),
					JournalID:         journalID,
					AccountID:         trxCategory.CreditAccountID,
					CompanyID:         input.CompanyID,
					Debit:             0,
					Credit:            amount,
					DebitAccountType:  trxCategory.DebitAccountType,
					CreditAccountType: trxCategory.CreditAccountType,
					TransactionType:   input.TransactionType,

					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		journalEntries = append(journalEntries, journal)
	}

	// Simpan semua journal entries ke database
	for _, entry := range journalEntries {
		if err := db.DB.Create(&entry).Error; err != nil {
			return nil, err
		}
	}

	return journalEntries, nil
}
