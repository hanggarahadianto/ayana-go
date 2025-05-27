package service

import (
	"ayana/db"
	"ayana/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateSingleJournalEntry(input models.JournalEntry) (models.JournalEntry, error) {
	var existing models.JournalEntry

	// Mulai transaksi
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Ambil data JournalEntry lama
		if err := tx.Preload("Lines").First(&existing, "id = ?", input.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("journal entry not found")
			}
			return err
		}

		// Hapus JournalLine lama
		if err := tx.Where("journal_id = ?", existing.ID).Delete(&models.JournalLine{}).Error; err != nil {
			return err
		}

		// Sinkronisasi dan siapkan JournalLine baru
		var newLines []models.JournalLine
		for _, line := range input.Lines {
			newLine := models.JournalLine{
				ID:                uuid.New(),
				JournalID:         existing.ID,
				AccountID:         line.AccountID,
				CompanyID:         input.CompanyID,
				Debit:             line.Debit,
				Credit:            line.Credit,
				Description:       input.Description, // disamakan
				TransactionType:   input.TransactionType,
				DebitAccountType:  input.DebitAccountType,
				CreditAccountType: input.CreditAccountType,
			}
			newLines = append(newLines, newLine)
		}

		// Update data JournalEntry
		existing.Transaction_ID = input.Transaction_ID
		existing.Invoice = input.Invoice
		existing.Description = input.Description
		existing.TransactionCategoryID = input.TransactionCategoryID
		existing.Amount = input.Amount
		existing.Partner = input.Partner
		existing.TransactionType = input.TransactionType
		existing.Status = input.Status
		existing.CompanyID = input.CompanyID
		existing.DateInputed = input.DateInputed
		existing.DueDate = input.DueDate
		existing.IsRepaid = input.IsRepaid
		existing.Installment = input.Installment
		existing.Note = input.Note
		existing.DebitAccountType = input.DebitAccountType
		existing.CreditAccountType = input.CreditAccountType

		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		// Simpan JournalLine baru
		if len(newLines) > 0 {
			if err := tx.Create(&newLines).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return models.JournalEntry{}, err
	}

	// Ambil ulang dengan preload
	if err := db.DB.Preload("Lines").First(&existing, "id = ?", existing.ID).Error; err != nil {
		return models.JournalEntry{}, err
	}

	return existing, nil
}
