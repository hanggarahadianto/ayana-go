// service/transaction_category_service.go

package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
)

type TransactionCategoryFilterParams struct {
	CompanyID       string            `json:"company_id"`
	TransactionType string            `json:"transaction_type"`
	Category        string            `json:"category"`
	Status          string            `json:"status"`
	Pagination      helper.Pagination `json:"pagination"`
}

func GetTransactionCategories(params TransactionCategoryFilterParams) ([]dto.TransactionCategoryResponse, int64, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
		Where("transaction_categories.company_id = ?", params.CompanyID)

	// Terapkan filter untuk transaction_type dan category
	tx = helper.ApplyTransactionFilters(tx, params.TransactionType, params.Category)

	if params.Status == "unpaid" {
		tx = tx.
			Where("transaction_categories.debit_account_type = ?", "Asset").
			Where("transaction_categories.credit_account_type = ?", "Liability")
	} else {
		// Batasi agar tidak mencatat transaksi antar Asset (contoh: kas ke bank dianggap tidak relevan di sini)
		tx = tx.Where(`
			NOT (
				transaction_categories.debit_account_type = 'Asset'
				AND transaction_categories.credit_account_type = 'Asset'
			)
		`).Where("transaction_categories.credit_account_type = ?", "Asset")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var transactions []models.TransactionCategory
	if params.Pagination.Limit == 0 {
		if err := tx.Preload("DebitAccount").Preload("CreditAccount").Find(&transactions).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := tx.Preload("DebitAccount").Preload("CreditAccount").
			Limit(params.Pagination.Limit).
			Offset(params.Pagination.Offset).
			Find(&transactions).Error; err != nil {
			return nil, 0, err
		}
	}

	responses := dto.MapToTransactionCategoryDTO(transactions)
	return responses, total, nil
}

func GetTransactionCategoriesWithoutPagination(companyID, transactionType, category string) ([]dto.TransactionCategoryResponse, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
		Where("transaction_categories.company_id = ?", companyID)

	tx = helper.ApplyTransactionFilters(tx, transactionType, category)

	var transactions []models.TransactionCategory
	if err := tx.Preload("DebitAccount").Preload("CreditAccount").Find(&transactions).Error; err != nil {
		return nil, err
	}

	responses := dto.MapToTransactionCategoryDTO(transactions)
	return responses, nil
}
