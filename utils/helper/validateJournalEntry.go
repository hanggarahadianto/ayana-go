package helper

import (
	"ayana/db"
	"ayana/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ValidateJournalEntry checks if the JournalEntry exists and transaction_id matches
func ValidateJournalEntry(id uuid.UUID, transactionID string) error {
	if id == uuid.Nil {
		return errors.New("ID is required")
	}

	if transactionID == "" {
		return errors.New("transaction ID not found")
	}

	// Fetch JournalEntry by ID
	var entry models.JournalEntry
	if err := db.DB.First(&entry, "id = ?", id).Error; err != nil {
		return fmt.Errorf("ID %s not found", id.String())
	}

	// Check if transaction_id matches exactly with database
	if entry.Transaction_ID != transactionID {
		return errors.New("transaction ID does not match with ID")
	}

	return nil
}
