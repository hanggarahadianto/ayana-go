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
		input.ID = uuid.New()
		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()

		if err := tx.Create(&input).Error; err != nil {
			return fmt.Errorf("failed to save to database: %w", err)
		}

		// Index ke Typesense
		if err := IndexCustomers(input); err != nil {
			return fmt.Errorf("indexing failed after DB commit: %w", err)
		}

		created = input
		return nil
	})

	if err != nil {
		return models.Customer{}, err
	}

	return created, nil
}
