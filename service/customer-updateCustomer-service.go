package service

import (
	"ayana/db"
	"ayana/models"
	"fmt"
)

func UpdateCustomerService(id string, input models.Customer) (*models.Customer, error) {
	var customer models.Customer

	// Cari customer lama
	if err := db.DB.First(&customer, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Update data customer
	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Status = input.Status
	customer.PaymentMethod = input.PaymentMethod
	customer.Amount = input.Amount
	customer.DateInputed = input.DateInputed
	customer.Marketer = input.Marketer
	customer.HomeID = input.HomeID
	customer.ProductUnit = input.ProductUnit
	customer.BankName = input.BankName

	// Simpan ke DB
	if err := db.DB.Save(&customer).Error; err != nil {
		return nil, fmt.Errorf("failed to save customer: %w", err)
	}

	// Update ke Typesense
	if err := updateCustomerInTypesense(customer); err != nil {
		return nil, fmt.Errorf("failed to update Typesense: %w", err)
	}

	return &customer, nil
}
