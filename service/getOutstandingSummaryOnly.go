package service

import (
	"ayana/db"
	"ayana/models"
)

// GetOutstandingDebtSummaryOnly menghitung total utang yang belum dibayar untuk perusahaan tertentu
func GetOutstandingDebtSummaryOnly(companyID string, status string) (int64, error) {
	var totalDebt int64

	// Subquery untuk filter hanya journal yang memiliki credit > 0 (hutang)
	subQuery := db.DB.
		Model(&models.JournalLine{}).
		Select("journal_id").
		Where("credit > 0 AND company_id = ?", companyID)

	// Hitung total utang (SUM amount) yang belum dibayar
	err := db.DB.Model(&models.JournalEntry{}).
		Where("id IN (?) AND status = ? AND is_repaid = false", subQuery, status).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalDebt).Error

	if err != nil {
		return 0, err // Return error jika query gagal
	}

	return totalDebt, nil
}
