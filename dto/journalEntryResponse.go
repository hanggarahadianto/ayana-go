package dto

import "time"

type JournalEntryResponse struct {
	ID                    string    `json:"id"`
	Invoice               string    `json:"invoice"`
	Description           string    `json:"description"`
	TransactionCategoryID string    `json:"transaction_category_id"`
	Amount                float64   `json:"amount"`
	Partner               string    `json:"partner"`
	TransactionType       string    `json:"transaction_type"`
	Status                string    `json:"status"`
	CompanyID             string    `json:"company_id"`
	DateInputed           time.Time `json:"date_inputed"`
	DueDate               time.Time `json:"due_date"`

	IsRepaid    bool                   `json:"is_repaid"`
	Installment int                    `json:"installment"`
	Note        string                 `json:"note"`
	Lines       []JournalEntryLineItem `json:"lines"`
}

type JournalEntryLineItem struct {
	ID             string  `json:"id"`
	JournalEntryID string  `json:"journal_entry_id"`
	AccountID      string  `json:"account_id"`
	AccountName    string  `json:"account_name"`
	Debit          float64 `json:"debit"`
	Credit         float64 `json:"credit"`
	Description    string  `json:"description"`
}
