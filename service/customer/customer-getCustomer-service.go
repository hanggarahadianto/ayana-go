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
	CompanyID   string
	Pagination  lib.Pagination
	DateFilter  lib.DateFilter
	SummaryOnly bool
	Search      string
	Status      string // ➕ Tambah ini
	SortBy      string
	SortOrder   string
}

func GetCustomersWithSearch(params CustomerFilterParams) ([]dto.CustomerResponse, int64, error) {
	var customers []models.Customer
	var total int64

	// 🔍 Search via Typesense
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
			return nil, int64(found), nil
		}
		return results, int64(found), nil
	}

	// Validasi sort_by dan sort_order
	validSortBy := map[string]bool{
		"date_inputed": true,
	}
	if !validSortBy[params.SortBy] {
		params.SortBy = "date_inputed"
	}
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "asc"
	}

	// 🔹 Query untuk menghitung total
	countQuery := db.DB.Model(&models.Customer{}).Where("company_id = ?", params.CompanyID)
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		countQuery = countQuery.Where("date_inputed BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		countQuery = countQuery.Where("date_inputed >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		countQuery = countQuery.Where("date_inputed <= ?", params.DateFilter.EndDate)
	}

	if params.Status != "" {
		countQuery = countQuery.Where("status = ?", params.Status)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 🔹 Query untuk ambil data (HARUS query baru)
	dataQuery := db.DB.Model(&models.Customer{}).
		Where("company_id = ?", params.CompanyID).
		Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Preload("Home").
		Preload("Marketer")

	// 🔹 Filter tanggal di dataQuery juga
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		dataQuery = dataQuery.Where("date_inputed BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		dataQuery = dataQuery.Where("date_inputed >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		dataQuery = dataQuery.Where("date_inputed <= ?", params.DateFilter.EndDate)
	}

	if params.Status != "" {
		dataQuery = dataQuery.Where("status = ?", params.Status)
	}

	if err := dataQuery.Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	// 🔹 Konversi ke DTO
	var response []dto.CustomerResponse
	for _, c := range customers {
		var home *dto.HomeResponse
		if c.Home.ID != uuid.Nil {
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
		response = append(response, dto.CustomerResponse{
			ID:            c.ID.String(),
			Name:          c.Name,
			Address:       c.Address,
			Phone:         c.Phone,
			Status:        c.Status,
			MarketerID:    c.MarketerID.String(),
			MarketerName:  c.MarketerName, // ✅ Langsung pakai dari field tabel
			Amount:        c.Amount,
			PaymentMethod: c.PaymentMethod,
			DateInputed:   c.DateInputed,
			HomeID:        c.HomeID.String(),
			ProductUnit:   c.ProductUnit,
			BankName:      c.BankName,
			Home:          home,
		})

	}

	return response, total, nil
}
