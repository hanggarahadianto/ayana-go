package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
)

// ✅ PARTIAL PAYMENT
func ProcessPartialPayment(existing models.JournalEntry, reversedEntry models.JournalEntry) ([]models.JournalEntry, error) {
	// 🔹 Hitung sisa hutang
	remainingAmount := existing.Amount - reversedEntry.Amount
	if remainingAmount < 0 {
		return nil, fmt.Errorf("payment amount cannot exceed existing amount")
	}
	fmt.Printf("exsiting amount: %d\n", existing.Amount)
	fmt.Printf("reversed amount: %d\n", reversedEntry.Amount)
	fmt.Println("remainingAmount", remainingAmount)

	// 🔹 Update existing → mark as unpaid & update amount
	if err := updateReserveJournalStatusPartial(existing.ID, remainingAmount); err != nil {
		return nil, fmt.Errorf("failed to update journal entry (partial payment): %w", err)
	}

	// 🔹 Ambil existing lengkap (dengan Lines)
	var existingWithLines models.JournalEntry
	if err := db.DB.Preload("Lines").First(&existingWithLines, "id = ?", existing.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch journal lines: %w", err)
	}

	// 🔹 Buat reversal entry untuk jumlah yang dibayar
	reversedEntry, err := createReserveJournalEntryService(existingWithLines, reversedEntry)
	if err != nil {
		return nil, err
	}

	return []models.JournalEntry{reversedEntry}, nil
}
