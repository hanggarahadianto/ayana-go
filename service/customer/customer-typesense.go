package service

import (
	"ayana/db"
	"ayana/dto"
	parse "ayana/lib"
	"ayana/models"
	tsClient "ayana/service"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/typesense/typesense-go/typesense/api"
)

func IndexCustomers(customers ...models.Customer) error {
	for _, customer := range customers {
		if err := IndexSingleCustomer(customer); err != nil {
			return fmt.Errorf("indexing failed for customer %s: %w", customer.ID, err)
		}
	}
	return nil
}

func IndexSingleCustomer(customer models.Customer) error {
	document := map[string]interface{}{
		"id":             customer.ID.String(),
		"name":           customer.Name,
		"address":        customer.Address,
		"phone":          customer.Phone,
		"status":         customer.Status,
		"amount":         customer.Amount,
		"payment_method": customer.PaymentMethod,
		"product_unit":   customer.ProductUnit,
		"bank_name":      customer.BankName,
		"company_id":     customer.CompanyID.String(),
	}

	if customer.DateInputed != nil {
		document["date_inputed"] = customer.DateInputed.Unix()
	}

	if customer.HomeID != nil {
		document["home_id"] = customer.HomeID.String()
	}

	// ✅ Tambahkan marketer_name (WAJIB karena ada di schema)
	document["marketer_name"] = customer.MarketerName

	// ✅ Tambahkan marketer_id jika ada
	if customer.MarketerID != uuid.Nil {
		document["marketer_id"] = customer.MarketerID.String()
	}

	_, err := tsClient.TsClient.Collection("customers").Documents().Create(context.Background(), document)
	if err != nil {
		return fmt.Errorf("gagal index customer ID %s: %w", customer.ID.String(), err)
	}
	return nil
}

func UpdateCustomerInTypesense(customer models.Customer) error {
	ctx := context.Background()

	docID := customer.ID.String()

	document := map[string]interface{}{
		"id":             docID,
		"name":           customer.Name,
		"address":        customer.Address,
		"phone":          customer.Phone,
		"status":         customer.Status,
		"payment_method": customer.PaymentMethod,
		"amount":         customer.Amount,
		"bank_name":      customer.BankName,
		"company_id":     customer.CompanyID.String(),
		"home_id":        customer.HomeID.String(),
		"product_unit":   customer.ProductUnit,
		"marketer_name":  customer.MarketerName,
		"marketer_id":    customer.MarketerID.String(),
		"date_inputed":   customer.DateInputed.Unix(),
	}

	// if customer.Marketer != nil {
	// 	document["marketer_name"] = customer.Marketer.Name
	// 	document["marketer_id"] = customer.Marketer.ID.String()
	// }
	// Langsung upsert
	_, err := tsClient.TsClient.Collection("customers").Documents().Upsert(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to upsert typesense document: %w", err)
	}

	return nil
}

func DeleteCustomerFromTypesense(ctx context.Context, customerIDs ...string) error {
	for _, id := range customerIDs {
		_, err := tsClient.TsClient.Collection("customers").Document(id).Delete(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "Not Found") {
				log.Printf("Typesense document not found for customer ID %s. Skipping.", id)
				continue
			}
			return fmt.Errorf("failed to delete customer %s: %w", id, err)
		}
	}
	return nil
}

func SearchCustomers(query, companyID string, startDate, endDate *time.Time, page, perPage int) ([]dto.CustomerResponse, int64, error) {
	log.Printf("🔍 Searching customers: query=%s, companyID=%s, page=%d, perPage=%d", query, companyID, page, perPage)

	// Build filter
	filters := []string{"company_id:=" + companyID}
	if startDate != nil {
		filters = append(filters, fmt.Sprintf("date_inputed:>=%d", startDate.Unix()))
	}
	if endDate != nil {
		filters = append(filters, fmt.Sprintf("date_inputed:<=%d", endDate.Unix()))
	}

	searchParams := &api.SearchCollectionParams{
		Q:        query,
		QueryBy:  "name,address,phone,status,marketer_name,bank_name",
		FilterBy: parse.PtrString(strings.Join(filters, " && ")),
		Page:     parse.PtrInt(page),
		PerPage:  parse.PtrInt(perPage),
	}

	searchResult, err := tsClient.TsClient.Collection("customers").Documents().Search(context.Background(), searchParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search in Typesense: %w", err)
	}
	log.Printf("🟢 Typesense search result:\nFound: %v\nSearchResult: %+v",
		searchResult.Found, searchResult)

	var results []dto.CustomerResponse
	var homeIDs []string
	var marketerIDs []string // 🔸 Tambah ini

	// Map search hits to response
	for _, hit := range *searchResult.Hits {
		doc := hit.Document
		if doc == nil {
			continue
		}
		m := *doc

		homeID := parse.GetString(m, "home_id")
		homeIDs = append(homeIDs, homeID)

		results = append(results, dto.CustomerResponse{
			ID:            parse.GetString(m, "id"),
			Name:          parse.GetString(m, "name"),
			Address:       parse.GetString(m, "address"),
			Phone:         parse.GetString(m, "phone"),
			Status:        parse.GetString(m, "status"),
			Amount:        parse.GetInt64(m, "amount"),
			PaymentMethod: parse.GetString(m, "payment_method"),
			DateInputed:   parse.GetTimePtr(m, "date_inputed"),
			HomeID:        homeID,
			ProductUnit:   parse.GetString(m, "product_unit"),
			BankName:      parse.GetString(m, "bank_name"),
			Marketer: &dto.MarketerResponse{
				ID:      parse.GetString(m, "marketer_id"),
				Name:    parse.GetString(m, "marketer_name"),
				IsAgent: parse.GetBool(m, "is_agent"),
			},
		})
	}

	// Ambil data home dari DB
	var homes []models.Home
	if len(homeIDs) > 0 {
		if err := db.DB.Where("id IN ?", homeIDs).Find(&homes).Error; err != nil {
			return nil, 0, fmt.Errorf("gagal mengambil data home: %w", err)
		}
	}
	// Ambil data marketer dari DB
	var marketers []models.Employee
	if len(marketerIDs) > 0 {
		if err := db.DB.Where("id IN ?", marketerIDs).Find(&marketers).Error; err != nil {
			return nil, 0, fmt.Errorf("gagal mengambil data marketer: %w", err)
		}
	}

	// Mapping home_id ke DTO
	homeMap := make(map[string]*dto.HomeResponse)
	for _, h := range homes {
		homeMap[h.ID.String()] = &dto.HomeResponse{
			ID:         h.ID.String(),
			ClusterID:  h.ClusterID.String(),
			Type:       h.Type,
			Title:      h.Title,
			Content:    h.Content,
			Bathroom:   int(h.Bathroom),
			Bedroom:    int(h.Bedroom),
			Square:     int(h.Square),
			Price:      int64(h.Price),
			Quantity:   int(h.Quantity),
			Status:     h.Status,
			Sequence:   int(h.Sequence),
			StartPrice: int64(h.StartPrice),
		}
	}

	for i, c := range results {
		if home, ok := homeMap[c.HomeID]; ok {
			results[i].Home = home
		}

	}

	var found int64
	if searchResult.Found != nil {
		found = int64(*searchResult.Found)
	}

	return results, found, nil
}

// 🔸 Helpers
