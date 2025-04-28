package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ExpenseFilterParams adalah parameter untuk filter pengeluaran
type AssetFilterParams struct {
	CompanyID       string
	Pagination      helper.Pagination
	DateFilter      helper.DateFilter
	AssetType       string
	TransactionType string
	SummaryOnly     bool
}

// Mengambil list dan summary pengeluaran
func GetAssets(params AssetFilterParams) ([]dto.JournalEntryResponse, int64, int64, error) {
	fmt.Printf("FILTER, params.AssetType: %s\n", params.AssetType)
	var entries []models.JournalEntry
	var total int64
	var totalAsset int64

	// Base query untuk semua jenis asset
	baseQuery := db.DB.Model(&models.JournalEntry{}).
		Where("debit_account_type = ?", "Asset").
		Where("company_id = ?", params.CompanyID)

	// Menambahkan filter berdasarkan jenis asset
	if params.AssetType == "fixed_asset" {
		baseQuery = baseQuery.Where("is_repaid = ? AND status = ? AND transaction_type = ?", true, "paid", "payout")
	} else if params.AssetType == "cashin" {
		baseQuery = baseQuery.Where("transaction_type = ?", "payin")
	} else if params.AssetType == "receivable" {
		baseQuery = baseQuery.Where("is_repaid = ? AND status = ? AND transaction_type = ?", false, "unpaid", "payin")
	}

	// Filter berdasarkan tanggal jika ada
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("date_inputed <= ?", params.DateFilter.EndDate)
	}

	// 1. Hitung total jumlah data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 2. Hitung total pengeluaran (tanpa ORDER BY, biar ga error)
	if err := baseQuery.Session(&gorm.Session{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalAsset).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jika summary_only = false, ambil data list aset
	if !params.SummaryOnly {
		dataQuery := baseQuery.Session(&gorm.Session{}).
			Order("date_inputed ASC").
			Limit(params.Pagination.Limit).
			Offset(params.Pagination.Offset)

		if err := dataQuery.Find(&entries).Error; err != nil {
			return nil, 0, 0, err
		}

		// Mapping ke DTO
		var responseData []dto.JournalEntryResponse
		for _, entry := range entries {
			var dueDate time.Time
			if entry.DueDate != nil {
				dueDate = *entry.DueDate // Menggunakan nilai sebenarnya
			} else {
				dueDate = time.Time{} // Nilai default untuk time.Time, yaitu zero value
			}

			responseData = append(responseData, dto.JournalEntryResponse{
				ID:                    entry.ID.String(),
				Invoice:               entry.Invoice,
				Description:           entry.Description,
				TransactionCategoryID: entry.TransactionCategoryID.String(),
				Amount:                float64(entry.Amount),
				Partner:               entry.Partner,
				TransactionType:       string(entry.TransactionType),
				Status:                string(entry.Status),
				CompanyID:             entry.CompanyID.String(),
				DateInputed:           *entry.DateInputed,
				DueDate:               dueDate, // Sekarang menggunakan tipe time.Time
				IsRepaid:              entry.IsRepaid,
				Installment:           entry.Installment,
				Note:                  entry.Note,
			})
		}

		return responseData, totalAsset, total, nil
	}

	// Jika summary_only = true, hanya kembalikan total dan totalAsset tanpa mengambil list data
	return nil, totalAsset, total, nil
}
