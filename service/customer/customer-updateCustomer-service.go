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

	// ✅ Ambil nama marketer baru
	var marketer models.Employee
	if err := db.DB.First(&marketer, "id = ?", input.MarketerID).Error; err != nil {
		return nil, fmt.Errorf("marketer not found: %w", err)
	}

	// Update field customer
	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Status = input.Status
	customer.PaymentMethod = input.PaymentMethod
	customer.Amount = input.Amount
	customer.DateInputed = input.DateInputed
	customer.MarketerID = input.MarketerID
	customer.MarketerName = marketer.Name // ✅ penting!
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
