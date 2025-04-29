package dto

import "time"

type JournalLineResponse struct {
	ID                string    `json:"id"`
	JournalEntryID    string    `json:"journal_entry_id"`
	Invoice           string    `json:"invoice"`
	Description       string    `json:"description"`
	Amount            float64   `json:"amount"`
	TransactionType   string    `json:"transaction_type"`
	DebitAccountType  string    `json:"debit_account_type"`
	CreditAccountType string    `json:"credit_account_type"`
	Status            string    `json:"status"`
	CompanyID         string    `json:"company_id"`
	DateInputed       time.Time `json:"date_inputed"`
	DueDate           time.Time `json:"due_date"`
	IsRepaid          bool      `json:"is_repaid"`
	Installment       int       `json:"installment"`
	Note              string    `json:"note"`
}
