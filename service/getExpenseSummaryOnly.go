package service

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
	"errors"

	"github.com/google/uuid"
)

func GetExpenseSummaryOnly(companyID string, dateFilter helper.DateFilter) (int64, error) {
	// Convert companyID ke UUID
	companyUUID, err := uuid.Parse(companyID)
	if err != nil {
		return 0, errors.New("format company_id tidak valid")
	}

	// Query untuk menghitung total pengeluaran
	var totalExpense int64
	query := db.DB.Model(&models.JournalEntry{}).
		Where("company_id = ?", companyUUID).
		Where("status = ?", "paid").
		Where("transaction_type = ?", "payout").
		Where("is_repaid = ?", true).
		Where("debit_account_type = ?", "Expense")

	// Filter berdasarkan tanggal
	if dateFilter.StartDate != nil {
		query = query.Where("date_inputed >= ?", dateFilter.StartDate)
	}
	if dateFilter.EndDate != nil {
		query = query.Where("date_inputed <= ?", dateFilter.EndDate)
	}

	err = query.Select("SUM(amount)").Scan(&totalExpense).Error
	if err != nil {
		return 0, err
	}

	return totalExpense, nil
}
