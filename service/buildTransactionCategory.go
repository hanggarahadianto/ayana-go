// service/transaction_category_service.go

package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"errors"
)

type TransactionCategoryFilterParams struct {
	CompanyID       string            `json:"company_id"`
	TransactionType string            `json:"transaction_type"`
	Category        string            `json:"category"`
	Status          string            `json:"status"`
	All             bool              `json:"all"` // ➕ baru
	Pagination      helper.Pagination `json:"pagination"`
}

func GetTransactionCategories(params TransactionCategoryFilterParams) ([]dto.TransactionCategoryResponse, int64, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
		Where("transaction_categories.company_id = ?", params.CompanyID)

	if !params.All {
		// ➕ Hanya gunakan filter jika All == false
		tx = helper.ApplyTransactionFilters(tx, params.TransactionType, params.Category, params.Status)

		// if params.Status == "unpaid" {
		// 	tx = tx.
		// 		Where("transaction_categories.debit_account_type = ?", "Asset").
		// 		Where("transaction_categories.credit_account_type = ?", "Liability")
		// } else {
		// 	tx = tx.Where(`
		// 		NOT (
		// 			transaction_categories.debit_account_type = 'Asset'
		// 			AND transaction_categories.credit_account_type = 'Asset'
		// 		)
		// 	`).Where("transaction_categories.debit_account_type = ?", "Asset").
		// 		Where("transaction_categories.credit_account_type = ?", "Asset")
		// }
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var transactions []models.TransactionCategory
	if err := tx.Preload("DebitAccount").Preload("CreditAccount").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	responses := dto.MapToTransactionCategoryDTO(transactions)
	return responses, total, nil
}

// func GetTransactionCategoriesWithoutPagination(companyID, transactionType, category string) ([]dto.TransactionCategoryResponse, error) {
// 	tx := db.DB.Model(&models.TransactionCategory{}).
// 		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
// 		Where("transaction_categories.company_id = ?", companyID)

// 	tx = helper.ApplyTransactionFilters(tx, transactionType, category)

// 	var transactions []models.TransactionCategory
// 	if err := tx.Preload("DebitAccount").Preload("CreditAccount").Find(&transactions).Error; err != nil {
// 		return nil, err
// 	}

// 	responses := dto.MapToTransactionCategoryDTO(transactions)
// 	return responses, nil
// }

func GetTransactionCategoriesWithoutPagination(params TransactionCategoryFilterParams) ([]dto.TransactionCategoryResponse, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id")

	// ✅ Pastikan minimal 1 filter aktif
	if params.CompanyID == "" && params.TransactionType == "" && params.Category == "" {
		return nil, errors.New("at least one filter (company_id, transaction_type, category) is required")
	}

	// ✅ Tambahkan filter sesuai yang dikirim
	if params.CompanyID != "" {
		tx = tx.Where("transaction_categories.company_id = ?", params.CompanyID)
	}
	if params.TransactionType != "" {
		tx = tx.Where("transaction_categories.transaction_type = ?", params.TransactionType)
	}
	if params.Category != "" {
		tx = tx.Where("transaction_categories.category = ?", params.Category)
	}

	var transactions []models.TransactionCategory
	if err := tx.Preload("DebitAccount").Preload("CreditAccount").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return dto.MapToTransactionCategoryDTO(transactions), nil
}
