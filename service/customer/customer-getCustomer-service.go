package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CustomerFilterParams struct {
	CompanyID    string
	Pagination   lib.Pagination
	DateFilter   lib.DateFilter
	SummaryOnly  bool
	Search       string
	Status       string // ➕ Tambah ini
	HasTestimony *bool  //
	SortBy       string
	SortOrder    string
}

func GetCustomersWithSearch(params CustomerFilterParams) ([]dto.CustomerResponse, int64, error) {
	var customers []models.Customer
	var total int64

	// 🔍 Gunakan Typesense jika ada search
	if params.Search != "" {
		results, found, err := SearchCustomers(
			params.Search,
			params.CompanyID,
			params.DateFilter.StartDate,
			params.DateFilter.EndDate,
			params.Pagination.Page,
			params.Pagination.Limit,
		)
		if err != nil {
			log.Println("🔴 Error search Typesense:", err)
			return nil, 0, fmt.Errorf("gagal search customer")
		}
		if params.SummaryOnly {
			return nil, found, nil
		}
		return results, found, nil
	}

	// ✅ Validasi sortBy & sortOrder
	if params.SortBy != "date_inputed" {
		params.SortBy = "date_inputed"
	}
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "asc"
	}

	// 🧱 Bangun base query untuk count dan fetch data
	baseQuery := db.DB.Model(&models.Customer{}).Where("customers.company_id = ?", params.CompanyID)

	if params.HasTestimony != nil {
		if *params.HasTestimony {
			baseQuery = baseQuery.
				Joins("JOIN testimonies ON testimonies.customer_id = customers.id").
				Where("testimonies.customer_id IS NOT NULL")
		} else {
			// LEFT JOIN dan IS NULL untuk customer yang tidak punya testimony
			baseQuery = baseQuery.
				Joins("LEFT JOIN testimonies ON testimonies.customer_id = customers.id").
				Where("testimonies.customer_id IS NULL")
		}
	}

	// 🔘 Filter tanggal
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("date_inputed <= ?", params.DateFilter.EndDate)
	}

	// 🔘 Filter status
	if params.Status != "" {
		baseQuery = baseQuery.Where("status = ?", params.Status)
	}

	// 🔢 Hitung total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 📦 Ambil data customer
	if err := baseQuery.
		Preload("Home").
		Preload("Marketer").
		Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	// 🔁 Mapping ke DTO
	response := make([]dto.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		var home *dto.HomeResponse
		if c.Home != nil && c.Home.ID != uuid.Nil {
			home = &dto.HomeResponse{
				ID:         c.Home.ID.String(),
				ClusterID:  c.Home.ClusterID.String(),
				Type:       c.Home.Type,
				Title:      c.Home.Title,
				Content:    c.Home.Content,
				Bathroom:   int(c.Home.Bathroom),
				Bedroom:    int(c.Home.Bedroom),
				Square:     int(c.Home.Square),
				Price:      int64(c.Home.Price),
				Quantity:   int(c.Home.Quantity),
				Status:     c.Home.Status,
				Sequence:   int(c.Home.Sequence),
				StartPrice: int64(c.Home.StartPrice),
			}
		}

		var marketer *dto.MarketerResponse
		if c.Marketer != nil {
			marketer = &dto.MarketerResponse{
				ID:      c.Marketer.ID.String(),
				Name:    c.Marketer.Name,
				IsAgent: c.Marketer.IsAgent,
			}
		} else {
			marketer = &dto.MarketerResponse{
				ID:      c.MarketerID.String(),
				Name:    "", // ✅ Hindari akses ke c.Marketer.Name karena nil
				IsAgent: false,
			}
		}

		response = append(response, dto.CustomerResponse{
			ID:            c.ID.String(),
			Name:          c.Name,
			Address:       c.Address,
			Phone:         c.Phone,
			Status:        c.Status,
			Amount:        c.Amount,
			PaymentMethod: c.PaymentMethod,
			DateInputed:   c.DateInputed,
			HomeID:        c.HomeID.String(),
			ProductUnit:   c.ProductUnit,
			BankName:      c.BankName,
			Home:          home,
			Marketer:      marketer,
		})
	}

	return response, total, nil
}
