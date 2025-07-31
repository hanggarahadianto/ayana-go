package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
)

type TransactionCategoryFilterParams struct {
	CompanyID         string
	TransactionType   string
	DebitCategory     string
	CreditCategory    string
	Status            string
	All               bool
	SelectOnly        bool
	DebitAccountType  string
	CreditAccountType string
	Pagination        lib.Pagination
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
	if params.DebitCategory != "" {
		tx = tx.Where("debit_category = ?", params.DebitCategory)
	}
	if params.CreditCategory != "" {
		tx = tx.Where("credit_category = ?", params.CreditCategory)
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
	if params.DebitAccountType != "" {
		tx = tx.Where("debit_account_type = ?", params.DebitAccountType)
	}
	if params.CreditAccountType != "" {
		tx = tx.Where("credit_account_type = ?", params.CreditAccountType)
	}
	if params.Status != "" {
		tx = tx.Where("status = ?", params.Status)
	}
	if params.DebitCategory != "" {
		tx = tx.Where("debit_category = ?", params.DebitCategory)
	}
	if params.CreditCategory != "" {
		tx = tx.Where("credit_category = ?", params.CreditCategory)
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
func GetUniqueCategories(params TransactionCategoryFilterParams) ([]string, string, error) {
	var categories []string

	if params.DebitAccountType == "" && params.CreditAccountType == "" {
		return nil, "Please provide either debit_account_type or credit_account_type", nil
	}

	tx := db.DB.Table("journal_entries AS je").
		Joins("JOIN transaction_categories AS tc ON tc.id = je.transaction_category_id").
		Where("je.company_id = ?", params.CompanyID).
		Where("je.transaction_type = ?", params.TransactionType)

	if params.Status != "" {
		tx = tx.Where("je.status = ?", params.Status)
	}
	if params.DebitAccountType != "" {
		tx = tx.Where("je.debit_account_type = ?", params.DebitAccountType)
	}
	if params.CreditAccountType != "" {
		tx = tx.Where("je.credit_account_type = ?", params.CreditAccountType)
	}

	var categoryColumn string
	if params.DebitAccountType != "" {
		categoryColumn = "credit_category" // ambil lawan dari debit
	} else {
		categoryColumn = "debit_category" // ambil lawan dari kredit
	}

	if err := tx.Distinct(categoryColumn).Pluck(categoryColumn, &categories).Error; err != nil {
		return nil, "", err
	}

	return categories, "success", nil
}
