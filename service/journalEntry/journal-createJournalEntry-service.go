package service

import (
	"ayana/db"
	"ayana/models"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func ProcessSingleJournalEntry(input models.JournalEntry) (models.JournalEntry, error) {
	return createJournalEntryService(input)
}

// Multiple entry (create & update)
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

func ProcessReverseJournalEntry(inputs []models.JournalEntry) ([]models.JournalEntry, error) {
	if len(inputs) == 0 {
		return nil, errors.New("input journal entries cannot be empty")
	}

	firstEntry := inputs[0]    // 10.000
	reversedEntry := inputs[1] /// 2.500

	fmt.Println("==== Reversed Entry ====")
	reversedJSON, _ := json.MarshalIndent(reversedEntry, "", "  ")
	fmt.Println(string(reversedJSON))

	// 🔹 Ambil existing journal entry
	var existing models.JournalEntry
	if err := db.DB.First(&existing, "id = ?", firstEntry.ID).Error; err != nil {
		return nil, fmt.Errorf("journal entry not found: %w", err)
	}

	switch {
	// case firstEntry.Amount == existing.Amount:
	// 	// FULL PAYMENT
	// 	results, err := ProcessFullPayment(existing, reversedEntry)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to process full payment: %w", err)
	// 	}
	// 	return results, nil

	case reversedEntry.Amount < existing.Amount:

		// PARTIAL PAYMENT
		results, err := ProcessPartialPayment(existing, reversedEntry)
		if err != nil {
			return nil, fmt.Errorf("failed to process full payment: %w", err)
		}
		return results, nil

	// TODO: tambahkan case partial / lainnya
	default:
		return nil, errors.New("unsupported payment scenario")
	}
}

// Core service: create / reverse journal
func createJournalEntryService(input models.JournalEntry) (models.JournalEntry, error) {
	if input.TransactionCategoryID == uuid.Nil || input.Amount <= 0 || input.CompanyID == uuid.Nil {
		return models.JournalEntry{}, fmt.Errorf("missing required fields")
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek company
	var company models.Company
	if err := tx.First(&company, "id = ?", input.CompanyID).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	// Ambil transaction category
	var trxCategory models.TransactionCategory
	if err := tx.Preload("DebitAccount").
		Preload("CreditAccount").
		First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	now := time.Now()
	journalID := uuid.New()

	// Build journal entry
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

	// Installment handling
	if input.Installment > 0 {
		installmentJournals, err := CreateInstallmentJournals(input)
		if err != nil {
			tx.Rollback()
			return models.JournalEntry{}, fmt.Errorf("create installment journal failed: %w", err)
		}
		if err := IndexJournals(append([]models.JournalEntry{journalWithDetails}, installmentJournals...)...); err != nil {
			return models.JournalEntry{}, fmt.Errorf("indexing failed after commit: %w", err)
		}
	} else {
		if err := IndexJournals(journalWithDetails); err != nil {
			return models.JournalEntry{}, fmt.Errorf("indexing failed after commit: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return models.JournalEntry{}, err
	}

	return journalWithDetails, nil
}

func createReserveJournalEntryService(existing models.JournalEntry, newInput models.JournalEntry) (models.JournalEntry, error) {
	fmt.Println("new input", newInput)
	if existing.ID == uuid.Nil || existing.CompanyID == uuid.Nil || len(existing.Lines) < 2 {
		return models.JournalEntry{}, fmt.Errorf("invalid journal entry for reversal")
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	journalID := uuid.New()

	// Ambil debit/credit dibalik dari existing
	debitAccountType := existing.CreditAccountType
	creditAccountType := existing.DebitAccountType
	debitAccountID := existing.Lines[1].AccountID  // credit lama jadi debit
	creditAccountID := existing.Lines[0].AccountID // debit lama jadi credit

	// Build reversed journal (mix antara existing dan newInput)
	reversed := models.JournalEntry{
		ID:                    journalID,
		Transaction_ID:        newInput.Transaction_ID, // ambil dari new input
		Invoice:               newInput.Invoice,        // ambil dari new input
		Description:           "Pelunasan untuk " + existing.Description,
		TransactionCategoryID: newInput.TransactionCategoryID, // tetap ikut existing
		Amount:                existing.Amount,
		DebitAccountType:      debitAccountType,
		CreditAccountType:     creditAccountType,
		Partner:               newInput.Partner,         // ambil dari new input
		TransactionType:       newInput.TransactionType, // ambil dari new input
		Status:                "done",
		IsRepaid:              true,
		DateInputed:           newInput.DateInputed,   // ambil dari new input
		RepaymentDate:         newInput.RepaymentDate, // ambil dari new input
		DueDate:               existing.DueDate,       // ikut existing
		Note:                  newInput.Note,          // ambil dari new input
		CompanyID:             existing.CompanyID,
		CreatedAt:             now,
		UpdatedAt:             now,
		Lines: []models.JournalLine{
			{
				ID:                uuid.New(),
				JournalID:         journalID,
				AccountID:         debitAccountID,
				CompanyID:         existing.CompanyID,
				Debit:             existing.Amount,
				Credit:            0,
				DebitAccountType:  debitAccountType,
				CreditAccountType: creditAccountType,
				TransactionType:   newInput.TransactionType, // ikut new input
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                uuid.New(),
				JournalID:         journalID,
				AccountID:         creditAccountID,
				CompanyID:         existing.CompanyID,
				Debit:             0,
				Credit:            existing.Amount,
				DebitAccountType:  debitAccountType,
				CreditAccountType: creditAccountType,
				TransactionType:   newInput.TransactionType, // ikut new input
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}

	if err := tx.Create(&reversed).Error; err != nil {
		tx.Rollback()
		return models.JournalEntry{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.JournalEntry{}, err
	}

	if err := IndexJournals(reversed); err != nil {
		return models.JournalEntry{}, fmt.Errorf("indexing failed after commit: %w", err)
	}

	return reversed, nil
}
