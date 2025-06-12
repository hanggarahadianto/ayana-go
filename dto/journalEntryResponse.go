package dto

import (
	"ayana/models"
	"time"
)

type JournalEntryResponse struct {
	ID                      string     `json:"id"`
	TransactionID           string     `json:"transaction_id"`
	TransactionCategoryID   string     `json:"transaction_category_id"`
	TransactionCategoryName string     `json:"transaction_category_name,omitempty"`
	Invoice                 string     `json:"invoice"`
	DebitCategory           string     `json:"debit_category,omitempty"`
	CreditCategory          string     `json:"credit_category,omitempty"`
	Partner                 string     `json:"partner"`
	Description             string     `json:"description"`
	Amount                  int64      `json:"amount"`
	TransactionType         string     `json:"transaction_type"`
	DebitAccountType        string     `json:"debit_account_type"`
	CreditAccountType       string     `json:"credit_account_type"`
	Status                  string     `json:"status"`
	CompanyID               string     `json:"company_id"`
	DateInputed             *time.Time `json:"date_inputed,omitempty"`
	DueDate                 *time.Time `json:"due_date,omitempty"`
	RepaymentDate           *time.Time `json:"repayment_date,omitempty"`
	IsRepaid                bool       `json:"is_repaid"`
	Installment             int        `json:"installment"`
	Note                    string     `json:"note"`
	PaymentDateStatus       string     `json:"payment_date_status,omitempty"`
	DebitLineId             string     `json:"debit_line_id,omitempty"`
	CreditLineId            string     `json:"credit_line_id,omitempty"`
	Label                   string     `json:"label,omitempty"`
}

func MapToJournalEntryResponse(entry models.JournalEntry) JournalEntryResponse {
	return JournalEntryResponse{
		ID:                      entry.ID.String(),
		TransactionID:           entry.Transaction_ID,
		TransactionCategoryID:   entry.TransactionCategoryID.String(),
		TransactionCategoryName: entry.TransactionCategory.Name, // pastikan relasi ini preload-ed
		Invoice:                 entry.Invoice,
		DebitCategory:           "", // Tidak tersedia langsung di model, bisa diisi dari relasi jika ada
		CreditCategory:          "", // Sama seperti atas
		Partner:                 entry.Partner,
		Description:             entry.Description,
		Amount:                  entry.Amount,
		TransactionType:         string(entry.TransactionType),
		DebitAccountType:        entry.DebitAccountType,
		CreditAccountType:       entry.CreditAccountType,
		Status:                  string(entry.Status),
		CompanyID:               entry.CompanyID.String(),
		DateInputed:             entry.DateInputed,
		DueDate:                 entry.DueDate,
		RepaymentDate:           entry.RepaymentDate,
		IsRepaid:                entry.IsRepaid,
		Installment:             entry.Installment,
		Note:                    entry.Note,
		PaymentDateStatus:       "", // Optional, tambahkan jika tersedia
		DebitLineId:             "", // Tambahkan jika ada di model/relasi
		CreditLineId:            "",
		Label:                   "", // Tambahkan jika ada
	}
}

func MapToJournalEntryResponses(entries []models.JournalEntry) []JournalEntryResponse {
	responses := make([]JournalEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = MapToJournalEntryResponse(entry)
	}
	return responses
}
