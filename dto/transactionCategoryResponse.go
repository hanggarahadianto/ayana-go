package dto

import "github.com/google/uuid"

type TransactionCategoryResponse struct {
	ID                uuid.UUID       `json:"id"`
	Name              string          `json:"name"`
	DebitAccountID    uuid.UUID       `json:"debit_account_id"`
	DebitAccountName  string          `json:"debit_account_name"`
	CreditAccountID   uuid.UUID       `json:"credit_account_id"`
	CreditAccountName string          `json:"credit_account_name"`
	Category          string          `json:"category"`
	Description       string          `json:"description"`
	CompanyID         uuid.UUID       `json:"company_id"`
	DebitAccount      AccountResponse `json:"debit_account"`
	CreditAccount     AccountResponse `json:"credit_account"`
}
