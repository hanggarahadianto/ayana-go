package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func updateReserveJournalStatusFullPayment(journalID uuid.UUID) error {
	result := db.DB.Model(&models.JournalEntry{}).
		Where("id = ?", journalID).
		Updates(map[string]interface{}{
			"is_repaid":  true,
			"status":     "done",
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("journal with ID %s not found", journalID)
	}

	return nil
}
func updateReserveJournalStatusPartial(journalID uuid.UUID, remainingAmount int64) error {
	result := db.DB.Model(&models.JournalEntry{}).
		Where("id = ?", journalID).
		Updates(map[string]interface{}{
			"is_repaid":  false,
			"status":     "unpaid",
			"amount":     remainingAmount, // ✅ update sisa hutang pakai int64
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("journal with ID %s not found", journalID)
	}

	return nil
}
