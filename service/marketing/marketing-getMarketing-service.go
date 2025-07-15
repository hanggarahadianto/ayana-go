package service

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	"sort"

	"gorm.io/gorm"
)

type MarketerFilterParams struct {
	CompanyID  string
	IsAgent    *bool // ➕ ditambahkan
	DateFilter lib.DateFilter
	SortBy     string
	SortOrder  string
}

type PerformanceMap map[string]*dto.PerformerResponse

func GetMarketerPerformance(params MarketerFilterParams) (dto.MarketerPerformanceResponse, error) {
	var employees []models.Employee

	query := db.DB.
		Where("company_id = ?", params.CompanyID).
		Where("LOWER(department) = ?", "marketing")

	if params.IsAgent != nil {
		query = query.Where("is_agent = ?", *params.IsAgent)
	}

	// Preload dengan filter tanggal ke relasi Customers
	query = query.Preload("Customers", func(db *gorm.DB) *gorm.DB {
		if params.DateFilter.StartDate != nil {
			db = db.Where("date_inputed >= ?", params.DateFilter.StartDate)
		}
		if params.DateFilter.EndDate != nil {
			db = db.Where("date_inputed <= ?", params.DateFilter.EndDate)
		}
		return db
	})

	if err := query.Find(&employees).Error; err != nil {
		return dto.MarketerPerformanceResponse{}, err
	}

	var performers []dto.PerformerResponse
	for _, emp := range employees {
		var totalAmount int64
		for _, cust := range emp.Customers {
			totalAmount += cust.Amount
		}
		performers = append(performers, dto.PerformerResponse{
			ID:           emp.ID.String(),
			Name:         emp.Name,
			TotalBooking: len(emp.Customers),
			TotalAmount:  totalAmount,
		})
	}

	// Sort
	sort.SliceStable(performers, func(i, j int) bool {
		return performers[i].TotalBooking > performers[j].TotalBooking
	})

	topN := 3
	if len(performers) < topN {
		topN = len(performers)
	}

	return dto.MarketerPerformanceResponse{
		TopPerformers:   performers[:topN],
		UnderPerformers: performers[topN:],
	}, nil
}
