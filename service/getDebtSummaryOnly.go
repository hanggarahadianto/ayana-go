package service

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
)

// Menghitung total outstanding debt saja
func GetOutstandingDebtSummaryOnly(companyID string, dateFilter helper.DateFilter) (int64, error) {
	var totalDebt int64

	subQuery := db.DB.
		Model(&models.JournalLine{}).
		Select("journal_id").
		Where("credit > 0 AND company_id = ?", companyID)

	query := db.DB.Model(&models.JournalEntry{}).
		Where("id IN (?) AND status = ? AND is_repaid = false", subQuery, "paid").
		Where("debit_account_type = ?", "Expense")

	// Tambahkan filter tanggal jika ada
	if dateFilter.StartDate != nil {
		query = query.Where("due_date >= ?", dateFilter.StartDate)
	}
	if dateFilter.EndDate != nil {
		query = query.Where("due_date <= ?", dateFilter.EndDate)
	}

	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&totalDebt).Error
	if err != nil {
		return 0, err
	}

	return totalDebt, nil
}
