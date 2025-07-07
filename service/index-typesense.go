package service

import (
	lib "ayana/lib"
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

var TsClient *typesense.Client

func InitTypesense(config *utilsEnv.Config) {
	TsClient = typesense.NewClient(
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

	journalIDs, err := fetchPostgresJournalIDs(db)
	if err != nil {
		return fmt.Errorf("gagal fetch journal IDs: %w", err)
	}
	if err := removeOrphanDocuments("journal_entries", journalIDs); err != nil {
		return fmt.Errorf("gagal hapus dokumen orphan journal: %w", err)
	}
	customerIDs, err := fetchPostgresCustomerIDs(db)
	if err != nil {
		return fmt.Errorf("gagal fetch customer IDs: %w", err)
	}
	if err := removeOrphanDocuments("customers", customerIDs); err != nil {
		return fmt.Errorf("gagal hapus dokumen orphan customer: %w", err)
	}

	log.Println("✅ Sinkronisasi selesai.")
	return nil
}
func boolPtr(b bool) *bool {
	return &b
}

func CreateCollectionIfNotExist() error {
	facetTrue := true
	defaultSort := "date_inputed"

	// === Journal Entries ===
	schema := &api.CollectionSchema{
		Name: "journal_entries",
		Fields: []api.Field{
			{Name: "id", Type: "string", Facet: &facetTrue},
			{Name: "company_id", Type: "string", Facet: &facetTrue},
			{Name: "debit_category", Type: "string", Facet: &facetTrue},
			{Name: "credit_category", Type: "string", Facet: &facetTrue},
			{Name: "transaction_id", Type: "string", Facet: &facetTrue},
			{Name: "invoice", Type: "string", Facet: &facetTrue},
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
			{Name: "status", Type: "string", Facet: &facetTrue},
			{Name: "installment", Type: "int32"},
			{Name: "note", Type: "string", Optional: boolPtr(true)},
			{Name: "payment_note", Type: "string", Optional: boolPtr(true)},
			{Name: "payment_note_color", Type: "string", Optional: boolPtr(true)},
		},
		DefaultSortingField: &defaultSort,
	}

	if _, err := TsClient.Collections().Create(context.Background(), schema); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println("⚠️  Collection 'journal_entries' sudah ada, lanjut...")
		} else {
			return fmt.Errorf("gagal buat collection journal_entries: %w", err)
		}
	} else {
		log.Println("✅ Collection 'journal_entries' berhasil dibuat")
	}

	// === Customers ===
	customerSchema := &api.CollectionSchema{
		Name: "customers",
		Fields: []api.Field{
			{Name: "id", Type: "string", Facet: &facetTrue},
			{Name: "name", Type: "string"},
			{Name: "address", Type: "string"},
			{Name: "phone", Type: "string"},
			{Name: "status", Type: "string", Facet: &facetTrue},
			{Name: "marketer_id", Type: "string", Facet: &facetTrue},
			{Name: "marketer_name", Type: "string"},
			{Name: "amount", Type: "float"},
			{Name: "payment_method", Type: "string"},
			{Name: "date_inputed", Type: "int64"}, // ⛔️ Jangan optional
			{Name: "home_id", Type: "string", Optional: boolPtr(true)},
			{Name: "product_unit", Type: "string"},
			{Name: "bank_name", Type: "string"},
			{Name: "company_id", Type: "string", Facet: &facetTrue},
		},
		DefaultSortingField: lib.PtrString("date_inputed"),
	}

	if _, err := TsClient.Collections().Create(context.Background(), customerSchema); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println("⚠️  Collection 'customers' sudah ada, lanjut...")
		} else {
			return fmt.Errorf("gagal buat collection customers: %w", err)
		}
	} else {
		log.Println("✅ Collection 'customers' berhasil dibuat")
	}

	return nil
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

func fetchPostgresCustomerIDs(db *gorm.DB) (map[string]struct{}, error) {
	var customers []models.Customer
	if err := db.Select("id").Find(&customers).Error; err != nil {
		return nil, err
	}

	ids := make(map[string]struct{}, len(customers))
	for _, c := range customers {
		ids[c.ID.String()] = struct{}{}
	}

	return ids, nil
}

func removeOrphanDocuments(collectionName string, validIDs map[string]struct{}) error {
	page := 1
	perPage := 250

	for {
		result, err := TsClient.Collection(collectionName).Documents().Search(context.Background(), &api.SearchCollectionParams{
			Q:       "*",
			Page:    &page,
			PerPage: &perPage,
		})
		if err != nil {
			return fmt.Errorf("search %s: %w", collectionName, err)
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

			if _, exists := validIDs[idStr]; exists {
				continue
			}

			if _, err := TsClient.Collection(collectionName).Document(idStr).Delete(context.Background()); err != nil {
				if strings.Contains(err.Error(), "404") {
					log.Printf("Dokumen %s (%s) tidak ditemukan, skip hapus", idStr, collectionName)
				} else {
					log.Printf("❌ Gagal hapus dokumen %s (%s): %v", idStr, collectionName, err)
				}
			} else {
				log.Printf("🗑️ Hapus dokumen %s dari collection %s", idStr, collectionName)
			}
		}

		if len(*hits) < perPage {
			break
		}
		page++
	}

	return nil
}
