// service/transaction_category_service.go

package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
)

type TransactionCategoryFilterParams struct {
	CompanyID       string `json:"company_id"`
	TransactionType string `json:"transaction_type"`
	Category        string `json:"category"`

	Pagination helper.Pagination `json:"pagination"`
}

// GetTransactionCategories mengambil data kategori transaksi berdasarkan parameter filter
func GetTransactionCategories(params TransactionCategoryFilterParams) ([]dto.TransactionCategoryResponse, int64, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
		Where("transaction_categories.company_id = ?", params.CompanyID)

	// ✅ Selalu terapkan semua filter
	tx = helper.ApplyTransactionFilters(tx, params.TransactionType, params.Category)

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
