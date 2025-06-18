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

func IndexCustomers(customers ...models.Customer) error {
	for _, customer := range customers {
		if err := IndexCustomerDocuments(customer); err != nil {
			return fmt.Errorf("indexing failed for customer %s: %w", customer.ID, err)
		}
	}
	return nil
}

func IndexCustomerDocuments(customers ...models.Customer) error {
	for _, customer := range customers {
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
		"marketer":       customer.Marketer,
		"amount":         customer.Amount,
		"payment_method": customer.PaymentMethod,
		"date_inputed": func() interface{} {
			if customer.DateInputed != nil {
				return customer.DateInputed.Unix()
			}
			return nil
		}(),
		"home_id": func() string {
			if customer.HomeID != nil {
				return customer.HomeID.String()
			}
			return ""
		}(),
		"product_unit": customer.ProductUnit,
		"bank_name":    customer.BankName,
		"created_at":   customer.CreatedAt.Unix(),
		"updated_at":   customer.UpdatedAt.Unix(),
	}

	_, err := tsClient.Collection("customers").Document(docID).Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	_, err = tsClient.Collection("customers").Document(docID).Update(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to update typesense document: %w", err)
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

func SearchCustomers(query, companyID string, page, perPage int) ([]dto.CustomerResponse, int64, error) {
	log.Printf("🔍 Searching customers: query=%s, companyID=%s, page=%d, perPage=%d", query, companyID, page, perPage)

	filters := []string{"company_id:=" + companyID}

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

		getInt64 := func(key string) int64 {
			if v, ok := m[key].(float64); ok {
				return int64(v)
			}
			return 0
		}

		getTimePtr := func(key string) *time.Time {
			if v, ok := m[key].(float64); ok && v > 0 {
				t := time.Unix(int64(v), 0)
				return &t
			}
			return nil
		}

		results = append(results, dto.CustomerResponse{
			ID:            getStr("id"),
			Name:          getStr("name"),
			Address:       getStr("address"),
			Phone:         getStr("phone"),
			Status:        getStr("status"),
			Marketer:      getStr("marketer"),
			Amount:        getInt64("amount"),
			PaymentMethod: getStr("payment_method"),
			DateInputed:   getTimePtr("date_inputed"),
			HomeID:        getStr("home_id"),
			ProductUnit:   getStr("product_unit"),
			BankName:      getStr("bank_name"),
		})
	}

	var found int64
	if searchResult.Found != nil {
		found = int64(*searchResult.Found)
	}

	return results, found, nil
}
