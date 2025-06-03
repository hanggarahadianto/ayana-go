package service

import (
	"ayana/dto"
	"ayana/models"
	utilsEnv "ayana/utils/env"
	"context"
	"log"
	"strings"
	"time"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

var tsClient *typesense.Client

func InitTypesense(config *utilsEnv.Config) {
	tsClient = typesense.NewClient(
		typesense.WithServer(config.TYPESENSE_HOST),
		typesense.WithAPIKey(config.TYPESENSE_API_KEY),
	)
}

func IndexJournalDocument(journal models.JournalEntry) error {
	document := map[string]interface{}{
		"id":                        journal.ID.String(),
		"transaction_id":            journal.Transaction_ID,
		"transaction_category_id":   journal.TransactionCategoryID.String(),
		"transaction_category_name": journal.TransactionCategory.Name,
		"invoice":                   journal.Invoice,
		"category":                  journal.TransactionCategory.Name,
		"partner":                   journal.Partner,
		"description":               journal.Description,
		"amount":                    journal.Amount,
		"transaction_type":          journal.TransactionType,
		"debit_account_type":        journal.DebitAccountType,
		"credit_account_type":       journal.CreditAccountType,
		"status":                    journal.Status,
		"company_id":                journal.CompanyID.String(),
		"date_inputed":              journal.DateInputed.Unix(),
		"due_date":                  journal.DueDate.Unix(),
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

func CreateCollectionIfNotExist() error {
	facetTrue := true
	defaultSort := "date_inputed"
	schema := &api.CollectionSchema{
		Name: "journal_entries",
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "company_id", Type: "string", Facet: &facetTrue},     // ✅ bisa di-facet
			{Name: "category", Type: "string", Facet: &facetTrue},       // ✅ bisa di-facet
			{Name: "transaction_id", Type: "string", Facet: &facetTrue}, // ✅ hanya satu kali
			{Name: "invoice", Type: "string", Facet: &facetTrue},        // ✅ bisa di-facet
			{Name: "partner", Type: "string"},
			{Name: "description", Type: "string"},
			{Name: "amount", Type: "float"},
			{Name: "date_inputed", Type: "int64", Facet: &facetTrue},
			{Name: "transaction_category_id", Type: "string"},
			{Name: "transaction_category_name", Type: "string"},
			{Name: "transaction_type", Type: "string"},
			{Name: "debit_account_type", Type: "string"},
			{Name: "credit_account_type", Type: "string"},
			{Name: "due_date", Type: "int64"},
			{Name: "repayment_date", Type: "int64"},
			{Name: "is_repaid", Type: "bool"},
			{Name: "installment", Type: "int32"},
			{Name: "note", Type: "string"},
		},
		DefaultSortingField: &defaultSort,
	}

	_, err := tsClient.Collections().Create(context.Background(), schema)
	if err != nil {
		// ✅ Jika collection sudah ada, abaikan error
		if strings.Contains(err.Error(), "already exists") {
			log.Println("⚠️  Collection 'journal_entries' sudah ada, lanjut...")
			return nil
		}
		// ❌ Jika error lain, baru return error
		return err
	}

	log.Println("✅ Collection 'journal_entries' berhasil dibuat")
	return nil
}

func SearchJournalLines(query string, companyID string, page, perPage int) ([]dto.JournalLineResponse, int, error) {
	searchParams := &api.SearchCollectionParams{
		Q: query,

		QueryBy:  "transaction_id,invoice,description,partner,category",
		FilterBy: ptrString("company_id:=" + companyID),
		Page:     ptrInt(page),
		PerPage:  ptrInt(perPage),
	}

	searchResult, err := tsClient.Collection("journal_entries").Documents().Search(context.Background(), searchParams)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.JournalLineResponse

	for _, hit := range *searchResult.Hits {
		doc := hit.Document
		if doc == nil {
			continue
		}
		m := *doc

		getStr := func(key string) string {
			if v, ok := m[key].(string); ok {
				return v
			}
			return ""
		}

		getFloat := func(key string) float64 {
			if v, ok := m[key].(float64); ok {
				return v
			}
			return 0
		}

		getInt := func(key string) int {
			if v, ok := m[key].(float64); ok {
				return int(v)
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

		getTime := func(key string) time.Time {
			if v, ok := m[key].(float64); ok {
				return time.Unix(int64(v), 0)
			}
			return time.Time{}
		}

		getTimePtr := func(key string) *time.Time {
			if v, ok := m[key].(float64); ok {
				t := time.Unix(int64(v), 0)
				return &t
			}
			return nil
		}

		results = append(results, dto.JournalLineResponse{
			ID:                      getStr("id"),
			Transaction_ID:          getStr("transaction_id"),
			TransactionCategoryID:   getStr("transaction_category_id"),
			TransactionCategoryName: getStr("transaction_category_name"),
			Invoice:                 getStr("invoice"),
			Category:                getStr("category"),
			Partner:                 getStr("partner"),
			Description:             getStr("description"),
			Amount:                  getFloat("amount"),
			TransactionType:         getStr("transaction_type"),
			DebitAccountType:        getStr("debit_account_type"),
			CreditAccountType:       getStr("credit_account_type"),
			Status:                  getStr("status"),
			CompanyID:               getStr("company_id"),
			DateInputed:             getTime("date_inputed"),
			DueDate:                 getTime("due_date"),
			RepaymentDate:           getTimePtr("repayment_date"),
			IsRepaid:                getBool("is_repaid"),
			Installment:             getInt("installment"),
			Note:                    getStr("note"),
			PaymentDateStatus:       getStr("payment_date_status"),
			DebitLineId:             getStr("debit_line_id"),
			CreditLineId:            getStr("credit_line_id"),
			Label:                   getStr("label"),
		})
	}

	found := 0
	if searchResult.Found != nil {
		found = *searchResult.Found
	}

	return results, found, nil
}

// Helper functions for pointers in api.SearchCollectionParams
func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
