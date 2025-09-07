package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
)

// ✅ FULL PAYMENT
// ✅ FULL PAYMENT
func ProcessFullPayment(existing models.JournalEntry, reversedEntry models.JournalEntry) ([]models.JournalEntry, error) {
	// 🔹 Update existing → mark as repaid & done
	if err := updateReserveJournalStatusFullPayment(existing.ID); err != nil {
		return nil, fmt.Errorf("failed to update journal entry (full payment): %w", err)
	}

	// 🔹 Ambil existing lengkap (dengan Lines)
	var existingWithLines models.JournalEntry
	if err := db.DB.Preload("Lines").First(&existingWithLines, "id = ?", existing.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch journal lines: %w", err)
	}

	// 🔹 Buat reversal entry berdasarkan existing (BUKAN payment)
	reversedEntry, err := createReserveJournalEntryService(existingWithLines, reversedEntry)
	if err != nil {
		return nil, err
	}

	// 🔹 Setelah reversal, baru simpan journal entry dari payment (normal create)

	return []models.JournalEntry{reversedEntry}, nil
}
