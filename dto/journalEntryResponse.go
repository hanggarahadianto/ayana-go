package dto

import (
	"ayana/models"
	"time"
)

type JournalEntryResponse struct {
	ID                    string    `json:"id"`
	Invoice               string    `json:"invoice"`
	TransactionID         string    `json:"transaction_id"`
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

// Fungsi untuk memetakan JournalEntry ke DTO
func MapToJournalEntryResponse(journalEntry models.JournalEntry) JournalEntryResponse {
	var lines []JournalEntryLineItem
	for _, line := range journalEntry.Lines {
		lines = append(lines, JournalEntryLineItem{
			ID:             line.ID.String(),
			JournalEntryID: line.JournalID.String(),
			AccountID:      line.AccountID.String(),
			AccountName:    line.Account.Name,
			Debit:          float64(line.Debit),
			Credit:         float64(line.Credit),
			Description:    line.Description,
		})
	}

	// Menangani nil untuk DateInputed dan DueDate
	var dateInputed, dueDate time.Time
	if journalEntry.DateInputed != nil {
		dateInputed = *journalEntry.DateInputed
	} else {
		// Menangani jika DateInputed nil
		dateInputed = time.Time{}
	}
	if journalEntry.DueDate != nil {
		dueDate = *journalEntry.DueDate
	} else {
		// Menangani jika DueDate nil
		dueDate = time.Time{}
	}

	return JournalEntryResponse{
		ID:                    journalEntry.ID.String(),
		Invoice:               journalEntry.Invoice,
		TransactionID:         journalEntry.Transaction_ID,
		Description:           journalEntry.Description,
		TransactionCategoryID: journalEntry.TransactionCategoryID.String(),
		Amount:                float64(journalEntry.Amount),
		Partner:               journalEntry.Partner,
		TransactionType:       string(journalEntry.TransactionType),
		Status:                string(journalEntry.Status),
		CompanyID:             journalEntry.CompanyID.String(),
		DateInputed:           dateInputed,
		DueDate:               dueDate,
		IsRepaid:              journalEntry.IsRepaid,
		Installment:           journalEntry.Installment,
		Note:                  journalEntry.Note,
		Lines:                 lines,
	}
}

// Memetakan banyak JournalEntries ke DTO
func MapToJournalEntryResponses(journalEntries []models.JournalEntry) []JournalEntryResponse {
	var responses []JournalEntryResponse
	for _, entry := range journalEntries {
		responses = append(responses, MapToJournalEntryResponse(entry))
	}
	return responses
}

func MapToJournalEntryResponseList(entries []models.JournalEntry) []JournalEntryResponse {
	var responseList []JournalEntryResponse
	for _, entry := range entries {
		responseList = append(responseList, MapToJournalEntryResponse(entry))
	}
	return responseList
}
