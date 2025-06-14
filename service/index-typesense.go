package service

import (
	"ayana/models"
	utilsEnv "ayana/utils/env"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"gorm.io/gorm"
)

var tsClient *typesense.Client

func InitTypesense(config *utilsEnv.Config) {
	tsClient = typesense.NewClient(
		typesense.WithServer(config.TYPESENSE_HOST),
		typesense.WithAPIKey(config.TYPESENSE_API_KEY),
	)
}

func SyncTypesenseWithPostgres(db *gorm.DB) error {
	log.Println("🔄 Sinkronisasi Typesense dengan PostgreSQL dimulai...")

	// Pastikan collection ada dulu
	if err := CreateCollectionIfNotExist(); err != nil {
		return fmt.Errorf("gagal pastikan collection: %w", err)
	}

	postgresIDs, err := fetchPostgresJournalIDs(db)
	if err != nil {
		return err
	}

	if err := removeOrphanDocuments(postgresIDs); err != nil {
		return err
	}

	log.Println("✅ Sinkronisasi selesai.")
	return nil
}

func CreateCollectionIfNotExist() error {
	facetTrue := true
	defaultSort := "date_inputed"
	schema := &api.CollectionSchema{
		Name: "journal_entries",
		Fields: []api.Field{
			{Name: "id", Type: "string", Facet: &facetTrue},
			{Name: "company_id", Type: "string", Facet: &facetTrue},      // ✅ bisa di-facet
			{Name: "debit_category", Type: "string", Facet: &facetTrue},  // ✅ bisa di-facet
			{Name: "credit_category", Type: "string", Facet: &facetTrue}, // ✅ bisa di-facet
			{Name: "transaction_id", Type: "string", Facet: &facetTrue},  // ✅ hanya satu kali
			{Name: "invoice", Type: "string", Facet: &facetTrue},         // ✅ bisa di-facet
			{Name: "partner", Type: "string"},
			{Name: "description", Type: "string"},
			{Name: "amount", Type: "float"},
			{Name: "date_inputed", Type: "int64", Facet: &facetTrue},
			{Name: "transaction_category_id", Type: "string"},
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

// Helper functions for pointers in api.SearchCollectionParams
func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}

func fetchPostgresJournalIDs(db *gorm.DB) (map[string]struct{}, error) {
	var journals []models.JournalEntry
	if err := db.Select("id").Find(&journals).Error; err != nil {
		return nil, err
	}

	ids := make(map[string]struct{}, len(journals))
	for _, j := range journals {
		ids[j.ID.String()] = struct{}{}
	}

	return ids, nil
}

func removeOrphanDocuments(validIDs map[string]struct{}) error {
	page := 1
	perPage := 250

	for {
		result, err := tsClient.Collection("journal_entries").Documents().Search(context.Background(), &api.SearchCollectionParams{
			Q:       "*",
			Page:    &page,
			PerPage: &perPage,
		})
		if err != nil {
			return err
		}

		hits := result.Hits
		if hits == nil || len(*hits) == 0 {
			break
		}

		for _, hit := range *hits {
			if hit.Document == nil {
				continue
			}
			doc := *hit.Document
			idVal, ok := doc["id"]
			if !ok {
				continue
			}
			idStr, ok := idVal.(string)
			if !ok {
				continue
			}

			// Cek apakah idStr ada di validIDs (data Postgres)
			if _, exists := validIDs[idStr]; exists {
				// Dokumen valid, tidak dihapus
				continue
			}

			// Jika tidak ada di validIDs, hapus dari Typesense
			_, err := tsClient.Collection("journal_entries").Document(idStr).Delete(context.Background())
			if err != nil {
				if strings.Contains(err.Error(), "404") {
					log.Printf("Dokumen %s tidak ditemukan, skip hapus", idStr)
				} else {
					return err
				}
			}
		}

		if len(*hits) < perPage {
			break
		}
		page++
	}

	return nil
}
