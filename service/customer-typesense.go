package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/typesense/typesense-go/typesense/api"
)

func IndexCustomers(customers ...models.Customer) error {
	for _, customer := range customers {
		if err := indexSingleCustomer(customer); err != nil {
			return fmt.Errorf("indexing failed for customer %s: %w", customer.ID, err)
		}
	}
	return nil
}

func indexSingleCustomer(customer models.Customer) error {
	document := map[string]interface{}{
		"id":             customer.ID.String(),
		"name":           customer.Name,
		"address":        customer.Address,
		"phone":          customer.Phone,
		"status":         customer.Status,
		"marketer":       customer.Marketer,
		"amount":         customer.Amount,
		"payment_method": customer.PaymentMethod,
		"product_unit":   customer.ProductUnit,
		"bank_name":      customer.BankName,
		"company_id":     customer.CompanyID.String(), // ✅ Tambahkan company_id
	}

	if customer.DateInputed != nil {
		document["date_inputed"] = customer.DateInputed.Unix()
	}

	if customer.HomeID != nil {
		document["home_id"] = customer.HomeID.String()
	}

	_, err := tsClient.Collection("customers").Documents().Create(context.Background(), document)
	if err != nil {
		return fmt.Errorf("gagal index customer ID %s: %w", customer.ID.String(), err)
	}
	return nil
}

func updateCustomerInTypesense(customer models.Customer) error {
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
		"marketer":       customer.Marketer,
		"date_inputed":   customer.DateInputed.Unix(),
	}

	// Langsung upsert
	_, err := tsClient.Collection("customers").Documents().Upsert(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to upsert typesense document: %w", err)
	}

	return nil
}

func DeleteCustomerFromTypesense(ctx context.Context, customerIDs ...string) error {
	for _, id := range customerIDs {
		_, err := tsClient.Collection("customers").Document(id).Delete(ctx)
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

	filters := []string{"company_id:=" + companyID}

	if startDate != nil {
		filters = append(filters, fmt.Sprintf("date_inputed:>=%d", startDate.Unix()))
	}
	if endDate != nil {
		filters = append(filters, fmt.Sprintf("date_inputed:<=%d", endDate.Unix()))
	}

	searchParams := &api.SearchCollectionParams{
		Q:        query,
		QueryBy:  "name,address,phone,status,marketer,bank_name",
		FilterBy: ptrString(strings.Join(filters, " && ")),
		Page:     ptrInt(page),
		PerPage:  ptrInt(perPage),
	}

	searchResult, err := tsClient.Collection("customers").Documents().Search(context.Background(), searchParams)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.CustomerResponse
	var homeIDs []string

	// Step 1: Parse basic customer fields from Typesense
	for _, hit := range *searchResult.Hits {
		doc := hit.Document
		if doc == nil {
			continue
		}
		m := *doc

		homeID := getStr(m, "home_id")
		homeIDs = append(homeIDs, homeID)

		results = append(results, dto.CustomerResponse{
			ID:            getStr(m, "id"),
			Name:          getStr(m, "name"),
			Address:       getStr(m, "address"),
			Phone:         getStr(m, "phone"),
			Status:        getStr(m, "status"),
			Marketer:      getStr(m, "marketer"),
			Amount:        getInt64(m, "amount"),
			PaymentMethod: getStr(m, "payment_method"),
			DateInputed:   getTimePtr(m, "date_inputed"),
			HomeID:        homeID,
			ProductUnit:   getStr(m, "product_unit"),
			BankName:      getStr(m, "bank_name"),
		})
	}

	// Step 2: Ambil data Home dari database
	var homes []models.Home
	if len(homeIDs) > 0 {
		err := db.DB.Where("id IN ?", homeIDs).Find(&homes).Error
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mengambil data home: %w", err)
		}
	}

	// Step 3: Map homeID → HomeResponse
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

	// Step 4: Lengkapi setiap CustomerResponse dengan Home
	for i, customer := range results {
		if home, ok := homeMap[customer.HomeID]; ok {
			results[i].Home = home
		}
	}

	// Total data ditemukan
	var found int64
	if searchResult.Found != nil {
		found = int64(*searchResult.Found)
	}

	return results, found, nil
}

// 🔸 Helpers

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func getTimePtr(m map[string]interface{}, key string) *time.Time {
	if v, ok := m[key].(float64); ok && v > 0 {
		t := time.Unix(int64(v), 0)
		return &t
	}
	return nil
}
