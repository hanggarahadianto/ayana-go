package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"time"
)

func UpdateSingleJournalEntry(input models.JournalEntry) (models.JournalEntry, error) {
	var existing models.JournalEntry

	// Cek apakah entry dengan ID tersebut ada
	if err := db.DB.First(&existing, "id = ?", input.ID).Error; err != nil {
		return models.JournalEntry{}, fmt.Errorf("journal entry not found")
	}

	// Lakukan update field yang diperlukan
	existing.Description = input.Note
	existing.TransactionCategoryID = input.TransactionCategoryID
	existing.Amount = input.Amount
	existing.Partner = input.Partner
	existing.TransactionType = input.TransactionType
	existing.Status = input.Status
	existing.IsRepaid = input.IsRepaid
	existing.DateInputed = input.DateInputed
	existing.DueDate = input.DueDate
	existing.Note = input.Note
	existing.UpdatedAt = time.Now()

	if err := db.DB.Save(&existing).Error; err != nil {
		return models.JournalEntry{}, err
	}

	// Reload dengan preload relasi
	var journalWithDetails models.JournalEntry
	if err := db.DB.Preload("Lines.Account").
		Preload("TransactionCategory.DebitAccount").
		Preload("TransactionCategory.CreditAccount").
		First(&journalWithDetails, "id = ?", input.ID).Error; err != nil {
		return models.JournalEntry{}, err
	}

	return journalWithDetails, nil
}
