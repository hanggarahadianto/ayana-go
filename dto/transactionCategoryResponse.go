// dto/transaction_category.go

package dto

import (
	"ayana/models"

	"github.com/google/uuid"
)

// TransactionCategoryResponse adalah response DTO untuk kategori transaksi
type TransactionCategoryResponse struct {
	ID                uuid.UUID       `json:"id"`
	Name              string          `json:"name"`
	DebitAccountID    uuid.UUID       `json:"debit_account_id"`
	DebitAccountType  string          `json:"debit_account_type"`
	CreditAccountID   uuid.UUID       `json:"credit_account_id"`
	CreditAccountType string          `json:"credit_account_type"`
	TransactionType   string          `json:"transaction_type"`
	Category          string          `json:"category"`
	Description       string          `json:"description"`
	CompanyID         uuid.UUID       `json:"company_id"`
	DebitAccount      AccountResponse `json:"debit_account"`
	CreditAccount     AccountResponse `json:"credit_account"`
}

// MapToTransactionCategoryDTO memetakan data transaksi ke DTO
func MapToTransactionCategoryDTO(transactions []models.TransactionCategory) []TransactionCategoryResponse {
	var responses []TransactionCategoryResponse
	for _, t := range transactions {
		res := TransactionCategoryResponse{
			ID:                t.ID,
			Name:              t.Name,
			DebitAccountID:    t.DebitAccountID,
			DebitAccountType:  t.DebitAccountType,
			CreditAccountID:   t.CreditAccountID,
			CreditAccountType: t.CreditAccountType,
			TransactionType:   t.TransactionType,
			Category:          t.Category,
			Description:       t.Description,
			CompanyID:         t.CompanyID,
			DebitAccount: AccountResponse{
				ID:          t.DebitAccount.ID,
				Code:        t.DebitAccount.Code,
				Name:        t.DebitAccount.Name,
				Type:        t.DebitAccount.Type,
				Category:    t.DebitAccount.Category,
				Description: t.DebitAccount.Description,
				CompanyID:   t.DebitAccount.CompanyID,
			},
			CreditAccount: AccountResponse{
				ID:          t.CreditAccount.ID,
				Code:        t.CreditAccount.Code,
				Name:        t.CreditAccount.Name,
				Type:        t.CreditAccount.Type,
				Category:    t.CreditAccount.Category,
				Description: t.CreditAccount.Description,
				CompanyID:   t.CreditAccount.CompanyID,
			},
		}
		responses = append(responses, res)
	}
	return responses
}
