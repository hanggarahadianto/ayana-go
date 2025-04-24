package service

import (
	"ayana/db"
	"database/sql"
)

// GetSummaryOnlyCashIn calculates the total cash in for a given company
func GetCashinSummaryOnly(companyID string) (int64, error) {
	var totalCashIn sql.NullInt64
	err := db.DB.Table("journal_entries").
		Where("company_id = ? AND transaction_type = ?", companyID, "payin").
		Select("SUM(amount)").
		Row().Scan(&totalCashIn)

	if err != nil {
		return 0, err // hanya jika query gagal, bukan karena NULL result
	}

	if totalCashIn.Valid {
		return totalCashIn.Int64, nil
	}

	return 0, nil // jika NULL (tidak ada data), kembalikan 0
}
