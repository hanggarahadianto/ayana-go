package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerFilterParams struct {
	CompanyID   string
	Pagination  helper.Pagination
	DateFilter  helper.DateFilter
	SummaryOnly bool
	Search      string
	SortBy      string
	SortOrder   string
}

func GetCustomersWithSearch(params CustomerFilterParams) ([]dto.CustomerResponse, int64, error) {
	var customers []models.Customer
	var total int64

	// Jika pakai search dengan Typesense (optional)
	if params.Search != "" {
		results, found, err := SearchCustomers(params.Search, params.CompanyID, params.Pagination.Page, params.Pagination.Limit)
		if err != nil {
			log.Println("🔴 Error search Typesense:", err)
			return nil, 0, fmt.Errorf("gagal search customer")
		}

		if params.SummaryOnly {
			return nil, int64(found), nil
		}

		return results, int64(found), nil
	}

	// 🔸 Inisialisasi baseQuery
	baseQuery := db.DB.Model(&models.Customer{})

	// 🔹 Filter tanggal jika ada
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("date_inputed BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("date_inputed >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("date_inputed <= ?", params.DateFilter.EndDate)
	}

	// 🔹 Hitung total data
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 🔹 Ambil data dengan pagination
	if err := baseQuery.
		Order("updated_at DESC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Preload("Home").
		Find(&customers).Error; err != nil {
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
				CreatedAt:  c.Home.CreatedAt,
				UpdatedAt:  c.Home.UpdatedAt,
			}
		}

		response = append(response, dto.CustomerResponse{
			ID:            c.ID.String(),
			Name:          c.Name,
			Address:       c.Address,
			Phone:         c.Phone,
			Status:        c.Status,
			Marketer:      c.Marketer,
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
