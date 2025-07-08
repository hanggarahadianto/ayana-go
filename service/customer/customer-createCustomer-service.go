package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCustomer(input models.Customer) (models.Customer, error) {
	var created models.Customer

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// ✅ Ambil nama marketer dari DB
		var marketer models.Employee
		if err := tx.First(&marketer, "id = ?", input.MarketerID).Error; err != nil {
			return fmt.Errorf("marketer not found: %w", err)
		}

		// Set UUID dan timestamps
		input.ID = uuid.New()
		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()
		input.MarketerName = marketer.Name

		// ✅ Simpan ke DB
		if err := tx.Create(&input).Error; err != nil {
			return fmt.Errorf("failed to save to database: %w", err)
		}

		// simpan hasil sementara
		created = input
		return nil
	})

	if err != nil {
		return models.Customer{}, err
	}

	// ✅ Index ke Typesense setelah DB sukses
	if err := IndexCustomers(created); err != nil {
		return models.Customer{}, fmt.Errorf("indexing failed: %w", err)
	}

	return created, nil
}
