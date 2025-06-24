package service

import (
	"ayana/db"
	"ayana/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func UpdateSingleJournalEntry(input models.JournalEntry) (models.JournalEntry, error) {
	var existing models.JournalEntry

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Ambil journal entry berdasarkan ID
		if err := tx.First(&existing, "id = ?", input.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("journal entry not found")
			}
			return err
		}

		// Update fields sesuai payload
		existing.Transaction_ID = input.Transaction_ID
		existing.Invoice = input.Invoice
		existing.Description = input.Note
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
		if input.DebitAccountType != "" {
			existing.DebitAccountType = input.DebitAccountType
		}

		if input.CreditAccountType != "" {
			existing.CreditAccountType = input.CreditAccountType
		}

		updates := map[string]interface{}{
			"transaction_type": input.TransactionType,
		}

		var lines []models.JournalLine
		if err := tx.Where("journal_id = ?", input.ID).Find(&lines).Error; err != nil {
			return err
		}

		for _, line := range lines {
			update := map[string]interface{}{
				"transaction_type": input.TransactionType,
			}

			// Debit dan Credit logic
			if line.Debit == 0 {
				update["debit"] = int64(0)
				update["credit"] = input.Amount
			} else if line.Credit == 0 {
				update["credit"] = int64(0)
				update["debit"] = input.Amount
			}

			// Hanya update jika value disediakan
			if input.DebitAccountType != "" {
				update["debit_account_type"] = input.DebitAccountType
			}
			if input.CreditAccountType != "" {
				update["credit_account_type"] = input.CreditAccountType
			}

			// Update per baris
			if err := tx.Model(&models.JournalLine{}).
				Where("id = ?", line.ID).
				Updates(update).Error; err != nil {
				return err
			}
		}

		// Eksekusi update ke JournalLine
		if err := tx.Model(&models.JournalLine{}).
			Where("journal_id = ?", input.ID).
			Updates(updates).Error; err != nil {
			return err
		}

		// Update JournalLine berdasarkan journal_id, hanya field yang diperlukan
		if len(updates) > 0 {
			if err := tx.Model(&models.JournalLine{}).
				Where("journal_id = ?", input.ID).
				Updates(updates).Error; err != nil {
				return err
			}
		}

		// Simpan perubahan journal entry
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.JournalEntry{}, err
	}
	err = UpdateJournalEntryInTypesense(existing)
	if err != nil {
		// Log error atau lakukan retry sesuai kebutuhan
		fmt.Printf("Warning: failed update to Typesense: %v\n", err)
	}

	return existing, nil
}
