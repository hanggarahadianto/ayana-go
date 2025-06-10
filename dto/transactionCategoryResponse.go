// dto/transaction_category.go

package dto

import (
	"ayana/models"

	"github.com/google/uuid"
)

// TransactionCategoryResponse adalah response DTO untuk kategori transaksi
type TransactionCategorySelectResponse struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Status            string    `json:"status"`
	TransactionType   string    `json:"transaction_type"`
	TransactionLabel  string    `json:"transaction_label"`
	DebitAccountType  string    `json:"debit_account_type"`
	CreditAccountType string    `json:"credit_account_type"`
	Description       string    `json:"description"`
}

func MapToTransactionCategorySelectDTO(data []models.TransactionCategory) []TransactionCategorySelectResponse {
	var result []TransactionCategorySelectResponse
	for _, item := range data {
		result = append(result, TransactionCategorySelectResponse{
			ID:                item.ID,
			Name:              item.Name,
			TransactionType:   item.TransactionType,
			Status:            item.Status,
			Description:       item.Description,
			TransactionLabel:  item.TransactionLabel, // hapus jika tidak perlu
			DebitAccountType:  item.DebitAccountType,
			CreditAccountType: item.CreditAccountType,
		})
	}
	return result
}

type TransactionCategoryResponse struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	TransactionType   string    `json:"transaction_type"`
	Status            string    `json:"status"`
	TransactionLabel  string    `json:"transaction_label"`
	DebitAccountID    uuid.UUID `json:"debit_account_id"`
	DebitAccountType  string    `json:"debit_account_type"`
	CreditAccountID   uuid.UUID `json:"credit_account_id"`
	CreditAccountType string    `json:"credit_account_type"`
	DebitCategory     string    `json:"debit_category"`
	CreditCategory    string    `json:"credit_category"`
	Description       string    `json:"description"`
	CompanyID         uuid.UUID `json:"company_id"`
}

// MapToTransactionCategoryDTO memetakan data transaksi ke DTO
func MapToTransactionCategoryDTO(transactions []models.TransactionCategory) []TransactionCategoryResponse {
	var responses []TransactionCategoryResponse
	for _, t := range transactions {

		res := TransactionCategoryResponse{
			ID:                t.ID,
			Name:              t.Name,
			Status:            t.Status,
			TransactionLabel:  t.TransactionLabel,
			DebitAccountID:    t.DebitAccountID,
			DebitAccountType:  t.DebitAccountType,
			CreditAccountID:   t.CreditAccountID,
			CreditAccountType: t.CreditAccountType,
			TransactionType:   t.TransactionType,
			DebitCategory:     t.DebitCategory,
			CreditCategory:    t.CreditCategory,
			Description:       t.Description,
			CompanyID:         t.CompanyID,
		}
		responses = append(responses, res)
	}
	return responses
}
