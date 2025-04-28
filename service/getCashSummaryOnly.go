package service

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
)

func GetCashSummaryOnly(companyID string, dateFilter helper.DateFilter) (int64, error) {
	var total int64

	err := db.DB.Model(&models.JournalEntry{}).
		Where("company_id = ? AND transaction_type = ?", companyID, "payin").
		Where("date_inputed BETWEEN ? AND ?", dateFilter.StartDate, dateFilter.EndDate).
		Select("SUM(amount)").Scan(&total).Error

	if err != nil {
		// Kalau error, tetap balikin 0
		return 0, nil
	}

	return total, nil
}
