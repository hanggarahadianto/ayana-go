package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"

	"gorm.io/gorm"
)

func UpdateCustomerService(id string, input models.Customer) (*models.Customer, error) {
	var updated models.Customer

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var customer models.Customer

		// Ambil data customer
		if err := tx.First(&customer, "id = ?", id).Error; err != nil {
			return fmt.Errorf("customer not found: %w", err)
		}

		// Ambil nama marketer
		var marketer models.Employee
		if err := tx.First(&marketer, "id = ?", input.MarketerID).Error; err != nil {
			return fmt.Errorf("marketer not found: %w", err)
		}

		// Update field
		customer.Name = input.Name
		customer.Address = input.Address
		customer.Phone = input.Phone
		customer.Status = input.Status
		customer.PaymentMethod = input.PaymentMethod
		customer.Amount = input.Amount
		customer.DateInputed = input.DateInputed
		customer.MarketerID = input.MarketerID
		customer.MarketerName = marketer.Name
		customer.HomeID = input.HomeID
		customer.ProductUnit = input.ProductUnit
		customer.BankName = input.BankName

		// Simpan ke DB
		if err := tx.Save(&customer).Error; err != nil {
			return fmt.Errorf("failed to save customer: %w", err)
		}

		// Simpan hasil ke variable luar
		updated = customer

		// ✅ Index ke Typesense
		if err := UpdateCustomerInTypesense(customer); err != nil {
			return fmt.Errorf("failed to update Typesense: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updated, nil
}
