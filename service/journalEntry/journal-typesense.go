package service

import (
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	tsClient "ayana/service"
	"ayana/utils/helper"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/typesense/typesense-go/typesense/api"
)

func IndexJournals(journals ...models.JournalEntry) error {
	for _, journal := range journals {
		if err := IndexJournalDocument(journal); err != nil {
			return fmt.Errorf("indexing failed for journal %s: %w", journal.ID, err)
		}
	}
	return nil
}

func IndexJournalDocument(journal models.JournalEntry) error {
	document := map[string]interface{}{
		"id":                      journal.ID.String(),
		"transaction_id":          journal.Transaction_ID,
		"transaction_category_id": journal.TransactionCategoryID.String(),
		"invoice":                 journal.Invoice,
		"debit_category":          journal.TransactionCategory.DebitCategory,
		"credit_category":         journal.TransactionCategory.CreditCategory,
		"partner":                 journal.Partner,
		"description":             journal.Description,
		"amount":                  journal.Amount,
		"transaction_type":        journal.TransactionType,
		"debit_account_type":      journal.DebitAccountType,
		"credit_account_type":     journal.CreditAccountType,
		"status":                  journal.Status,
		"company_id":              journal.CompanyID.String(),
		"date_inputed":            journal.DateInputed.Unix(),
		"due_date":                journal.DueDate.Unix(),
		"repayment_date": func() interface{} {
			if journal.RepaymentDate != nil {
				return journal.RepaymentDate.Unix()
			}
			return nil
		}(),
		"is_repaid":   journal.IsRepaid,
		"installment": journal.Installment,
		"note":        journal.Note,
	}

	_, err := tsClient.TsClient.Collection("journal_entries").Documents().Create(context.Background(), document)
	return err
}

func UpdateJournalEntryInTypesense(entry models.JournalEntry) error {
	ctx := context.Background() // Buat context

	docID := entry.ID.String()

	document := map[string]interface{}{
		"id":                      docID,
		"transaction_id":          entry.Transaction_ID,
		"invoice":                 entry.Invoice,
		"description":             entry.Description,
		"transaction_category_id": entry.TransactionCategoryID.String(),
		"amount":                  entry.Amount,
		"partner":                 entry.Partner,
		"transaction_type":        entry.TransactionType,
		"status":                  entry.Status,
		"company_id":              entry.CompanyID.String(),
		"date_inputed":            entry.DateInputed.Unix(),
		"due_date":                entry.DueDate.Unix(),
		"is_repaid":               entry.IsRepaid,
		"installment":             entry.Installment,
		"note":                    entry.Note,
		"debit_account_type":      entry.DebitAccountType,
		"credit_account_type":     entry.CreditAccountType,
	}

	// Cek dokumen sudah ada atau belum
	_, err := tsClient.TsClient.Collection("journal_entries").Document(docID).Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("document not found in typesense: %w", err)
	}

	// Update dokumen
	_, err = tsClient.TsClient.Collection("journal_entries").Document(docID).Update(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to update typesense document: %w", err)
	}

	return nil
}

func DeleteJournalEntryFromTypesense(ctx context.Context, journalEntryIDs ...string) error {
	for _, id := range journalEntryIDs {
		_, err := tsClient.TsClient.
			Collection("journal_entries").
			Document(id).
			Delete(ctx)

		if err != nil {
			// Abaikan error "Not Found"
			if strings.Contains(err.Error(), "Not Found") {
				log.Printf("Typesense document not found for ID %s. Skipping deletion.", id)
				continue
			}

			// Log error dan kembalikan
			log.Printf("Failed to delete document %s: %v", id, err)
			return fmt.Errorf("failed to delete document %s from Typesense: %w", id, err)
		}
	}

	return nil
}

func SearchJournalLines(
	query string,
	companyID string,
	accountType string,
	debitCategory string,
	creditCategory string,
	startDate, endDate *time.Time,
	Type *string,
	page, perPage int,
) ([]dto.JournalEntryResponse, int, error) {
	log.Printf("🔍 Searching journal lines: query=%s, companyID=%s,AccountType=%v,Type=%v,page=%d, perPage=%d", query, companyID, accountType, Type, page, perPage)

	// ✅ Gunakan lib refactor untuk filter
	filterBy := helper.BuildTypesenseFilter(companyID, accountType, debitCategory, creditCategory, startDate, endDate, Type)

	// 🔧 Setup search params
	searchParams := &api.SearchCollectionParams{
		Q:        query,
		QueryBy:  "transaction_id,invoice,description,partner,debit_category,credit_category",
		FilterBy: &filterBy,
		Page:     &page,
		PerPage:  &perPage,
		SortBy:   lib.PtrString("date_inputed:desc"),
	}

	// 🧪 Debug log
	log.Println("🔥 Typesense Search Params:")
	log.Printf("Q         : %s\n", searchParams.Q)
	log.Printf("QueryBy   : %s\n", searchParams.QueryBy)
	log.Printf("FilterBy  : %s\n", *searchParams.FilterBy)
	log.Printf("Page      : %d\n", *searchParams.Page)
	log.Printf("PerPage   : %d\n", *searchParams.PerPage)
	log.Printf("SortBy    : %s\n", *searchParams.SortBy)

	// 🚀 Perform search
	searchResult, err := tsClient.TsClient.Collection("journal_entries").Documents().Search(context.Background(), searchParams)
	if err != nil {
		return nil, 0, err
	}

	// 🧾 Parse results
	var results []dto.JournalEntryResponse
	for _, hit := range *searchResult.Hits {
		doc := hit.Document
		if doc == nil {
			continue
		}
		m := *doc

		results = append(results, dto.JournalEntryResponse{
			ID:                      lib.GetString(m, "id"),
			TransactionID:           lib.GetString(m, "transaction_id"),
			TransactionCategoryID:   lib.GetString(m, "transaction_category_id"),
			TransactionCategoryName: lib.GetString(m, "transaction_category_name"),
			Invoice:                 lib.GetString(m, "invoice"),
			DebitCategory:           lib.GetString(m, "debit_category"),
			CreditCategory:          lib.GetString(m, "credit_category"),
			Partner:                 lib.GetString(m, "partner"),
			Description:             lib.GetString(m, "description"),
			Amount:                  lib.GetInt64(m, "amount"),
			TransactionType:         lib.GetString(m, "transaction_type"),
			DebitAccountType:        lib.GetString(m, "debit_account_type"),
			CreditAccountType:       lib.GetString(m, "credit_account_type"),
			Status:                  lib.GetString(m, "status"),
			CompanyID:               lib.GetString(m, "company_id"),
			DateInputed:             lib.GetTimePtr(m, "date_inputed"),
			DueDate:                 lib.GetTimePtr(m, "due_date"),
			RepaymentDate:           lib.GetTimePtr(m, "repayment_date"),
			IsRepaid:                lib.GetBool(m, "is_repaid"),
			Installment:             lib.GetInt(m, "installment"),
			Note:                    lib.GetString(m, "note"),
			PaymentNote:             lib.GetString(m, "payment_note"),
			PaymentNoteColor:        lib.GetString(m, "payment_note_color"),
			DebitLineId:             lib.GetString(m, "debit_line_id"),
			CreditLineId:            lib.GetString(m, "credit_line_id"),
			Label:                   lib.GetString(m, "label"),
		})
	}

	found := 0
	if searchResult.Found != nil {
		found = *searchResult.Found
	}
	return results, found, nil
}
