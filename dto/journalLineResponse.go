package dto

import "time"

type JournalLineResponse struct {
	ID                      string    `json:"id"`
	Transaction_ID          string    `json:"transaction_id"`
	TransactionCategoryID   string    `json:"transaction_category_id"`
	TransactionCategoryName string    `json:"transaction_category_name,omitempty"` // optional, jika perlu
	Invoice                 string    `json:"invoice"`
	Category                string    `json:"category"`
	Partner                 string    `json:"partner"`
	Description             string    `json:"description"`
	Amount                  float64   `json:"amount"`
	TransactionType         string    `json:"transaction_type"`
	DebitAccountType        string    `json:"debit_account_type"`
	CreditAccountType       string    `json:"credit_account_type"`
	Status                  string    `json:"status"`
	CompanyID               string    `json:"company_id"`
	DateInputed             time.Time `json:"date_inputed"`
	DueDate                 time.Time `json:"due_date"`
	IsRepaid                bool      `json:"is_repaid"`
	Installment             int       `json:"installment"`
	Note                    string    `json:"note"`
	PaymentDateStatus       string    `json:"payment_date_status,omitempty"`
	DebitLineId             string    `json:"debit_line_id,omitempty"`  // optional, jika perlu
	CreditLineId            string    `json:"credit_line_id,omitempty"` // optional, jika perlu
	Label                   string    `json:"label,omitempty"`          // optional, jika perlu
}
