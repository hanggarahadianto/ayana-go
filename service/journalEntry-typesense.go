package service

import (
	"ayana/dto"
	"ayana/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/typesense/typesense-go/typesense/api"
)

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

	_, err := tsClient.Collection("journal_entries").Documents().Create(context.Background(), document)
	return err
}

func updateJournalEntryInTypesense(entry models.JournalEntry) error {
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
	_, err := tsClient.Collection("journal_entries").Document(docID).Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("document not found in typesense: %w", err)
	}

	// Update dokumen
	_, err = tsClient.Collection("journal_entries").Document(docID).Update(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to update typesense document: %w", err)
	}

	return nil
}

func DeleteJournalEntryFromTypesense(ctx context.Context, journalEntryID string) error {
	_, err := tsClient.
		Collection("journal_entries").
		Document(journalEntryID).
		Delete(ctx)

	if err != nil {
		// Jika error karena dokumen tidak ditemukan, abaikan
		if strings.Contains(err.Error(), "Not Found") {
			log.Printf("Typesense document not found for ID %s. Skipping deletion.", journalEntryID)
			return nil
		}
		// Untuk error lain, tetap return error
		return fmt.Errorf("failed to delete document from Typesense: %w", err)
	}

	return nil
}

func SearchJournalLines(query string, companyID string, debitCategory string, creditCategory string, page, perPage int) ([]dto.JournalEntryResponse, int, error) {
	log.Printf("🔍 Searching journal lines: query=%s, companyID=%s, page=%d, perPage=%d", query, companyID, page, perPage)

	filters := []string{"company_id:=" + companyID}
	if debitCategory != "" {
		filters = append(filters, fmt.Sprintf("debit_category:=%q", debitCategory))
	}
	if creditCategory != "" {
		filters = append(filters, fmt.Sprintf("credit_category:=%q", creditCategory))
	}

	searchParams := &api.SearchCollectionParams{
		Q:        query,
		QueryBy:  "transaction_id,invoice,description,partner,debit_category,credit_category",
		FilterBy: ptrString(strings.Join(filters, " && ")),
		Page:     ptrInt(page),
		PerPage:  ptrInt(perPage),
	}

	searchResult, err := tsClient.Collection("journal_entries").Documents().Search(context.Background(), searchParams)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.JournalEntryResponse

	for _, hit := range *searchResult.Hits {
		doc := hit.Document
		if doc == nil {
			continue
		}
		m := *doc

		get := func(key string) string {
			if v, ok := m[key].(string); ok {
				return v
			}
			return ""
		}
		getInt := func(key string) int {
			if v, ok := m[key].(float64); ok {
				return int(v)
			}
			return 0
		}
		getInt64 := func(key string) int64 {
			if v, ok := m[key].(float64); ok {
				return int64(v)
			}
			return 0
		}
		getBool := func(key string) bool {
			switch v := m[key].(type) {
			case bool:
				return v
			case string:
				return v == "true"
			}
			return false
		}
		getTimePtr := func(key string) *time.Time {
			if v, ok := m[key].(float64); ok && v > 0 {
				t := time.Unix(int64(v), 0)
				return &t
			}
			return nil
		}

		results = append(results, dto.JournalEntryResponse{
			ID:                      get("id"),
			TransactionID:           get("transaction_id"),
			TransactionCategoryID:   get("transaction_category_id"),
			TransactionCategoryName: get("transaction_category_name"),
			Invoice:                 get("invoice"),
			DebitCategory:           get("debit_category"),
			CreditCategory:          get("credit_category"),
			Partner:                 get("partner"),
			Description:             get("description"),
			Amount:                  getInt64("amount"),
			TransactionType:         get("transaction_type"),
			DebitAccountType:        get("debit_account_type"),
			CreditAccountType:       get("credit_account_type"),
			Status:                  get("status"),
			CompanyID:               get("company_id"),
			DateInputed:             getTimePtr("date_inputed"),
			DueDate:                 getTimePtr("due_date"),
			RepaymentDate:           getTimePtr("repayment_date"),
			IsRepaid:                getBool("is_repaid"),
			Installment:             getInt("installment"),
			Note:                    get("note"),
			PaymentDateStatus:       get("payment_date_status"),
			DebitLineId:             get("debit_line_id"),
			CreditLineId:            get("credit_line_id"),
			Label:                   get("label"),
		})
	}

	found := 0
	if searchResult.Found != nil {
		found = *searchResult.Found
	}

	return results, found, nil
}
