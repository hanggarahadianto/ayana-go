// service/transaction_category_service.go

package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
)

type TransactionCategoryFilterParams struct {
	CompanyID       string
	TransactionType string
	Category        string
	Status          string
	All             bool
	SelectOnly      bool
	Pagination      helper.Pagination
}

func GetTransactionCategoriesAll() ([]dto.TransactionCategoryResponse, error) {
	var categories []models.TransactionCategory
	err := db.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return dto.MapToTransactionCategoryDTO(categories), nil
}

func GetTransactionCategoriesForSelect(params TransactionCategoryFilterParams) ([]dto.TransactionCategorySelectResponse, error) {
	tx := db.DB.Model(&models.TransactionCategory{})

	// Wajib filter CompanyID & optional filter lainnya
	tx = tx.Where("company_id = ?", params.CompanyID)

	if params.TransactionType != "" {
		tx = tx.Where("transaction_type = ?", params.TransactionType)
	}
	if params.Status != "" {
		tx = tx.Where("status = ?", params.Status)
	}
	if params.Category != "" {
		tx = tx.Where("category = ?", params.Category)
	}

	var categories []models.TransactionCategory
	if err := tx.Find(&categories).Error; err != nil {
		return nil, err
	}
	return dto.MapToTransactionCategorySelectDTO(categories), nil
}

func GetTransactionCategoriesWithPagination(params TransactionCategoryFilterParams) ([]dto.TransactionCategoryResponse, int64, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Where("company_id = ?", params.CompanyID)

	if params.TransactionType != "" {
		tx = tx.Where("transaction_type = ?", params.TransactionType)
	}
	if params.Status != "" {
		tx = tx.Where("status = ?", params.Status)
	}
	if params.Category != "" {
		tx = tx.Where("category = ?", params.Category)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var categories []models.TransactionCategory
	if err := tx.Order("updated_at DESC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return dto.MapToTransactionCategoryDTO(categories), total, nil
}

func GetUniqueCategories(params TransactionCategoryFilterParams) ([]string, error) {
	tx := db.DB.Model(&models.TransactionCategory{}).
		Select("DISTINCT category").
		Where("company_id = ?", params.CompanyID)

	if params.TransactionType != "" {
		tx = tx.Where("transaction_type = ?", params.TransactionType)
	}
	if params.Status != "" {
		tx = tx.Where("status = ?", params.Status)
	}

	var categories []string
	if err := tx.Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
