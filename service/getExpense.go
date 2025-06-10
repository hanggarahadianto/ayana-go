package service

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"fmt"
	"log"
	"math"

	"gorm.io/gorm"
)

type ExpenseFilterParams struct {
	CompanyID     string
	Pagination    helper.Pagination
	DateFilter    helper.DateFilter
	SummaryOnly   bool
	ExpenseStatus string
	Category      string
	Search        string // ⬅️ Tambahkan ini
}

func GetExpensesFromJournalLines(params ExpenseFilterParams) ([]dto.JournalLineResponse, int64, int64, error) {
	var lines []models.JournalLine
	var total int64
	var totalExpense int64

	if params.Search != "" {
		results, found, err := SearchJournalLines(params.Search, params.CompanyID, params.Category, params.Pagination.Page, params.Pagination.Limit)

		if err != nil {
			log.Println("Error saat search ke Typesense:", err)
			return nil, 0, 0, fmt.Errorf("gagal mengambil data aset: %w", err)
		}

		// Jika hanya summary diperlukan
		if params.SummaryOnly {
			var totalExpense int64 = 0
			for _, line := range results {
				totalExpense += int64(line.Amount)
			}
			return nil, totalExpense, int64(found), nil
		}

		var totalExpense int64 = 0
		for _, line := range results {
			totalExpense += int64(line.Amount)
		}

		return results, totalExpense, int64(found), nil
	}

	// Build base query tanpa Limit dan Offset untuk Count dan Sum
	baseQuery := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Joins("LEFT JOIN transaction_categories ON journal_entries.transaction_category_id = transaction_categories.id").
		Where("journal_entries.company_id = ?", params.CompanyID).
		Where("journal_lines.debit_account_type = ?", "Expense")

	// Filter expense type
	switch params.ExpenseStatus {
	case "base":
		baseQuery = baseQuery.
			Where("journal_lines.debit > 0"). // Expense biasanya di sisi debit
			Where("journal_entries.is_repaid = ? AND journal_entries.status = ? AND journal_entries.transaction_type = ?", true, "paid", "payout")
	}

	if params.Category != "" {
		baseQuery = baseQuery.Where("transaction_categories.category = ?", params.Category)
	}

	// Filter date
	if params.DateFilter.StartDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed >= ?", params.DateFilter.StartDate)
	}
	if params.DateFilter.EndDate != nil {
		baseQuery = baseQuery.Where("journal_entries.date_inputed <= ?", params.DateFilter.EndDate)
	}

	// Hitung total baris
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Hitung total expense (menggunakan debit - credit)
	if err := baseQuery.Session(&gorm.Session{}).
		Select("COALESCE(SUM(journal_lines.debit - journal_lines.credit), 0)").
		Scan(&totalExpense).Error; err != nil {
		return nil, 0, 0, err
	}

	// // Jika SummaryOnly = true, kembalikan hanya totalExpense dan total
	// if params.SummaryOnly {
	// 	return nil, totalExpense, total, nil
	// }

	// Query untuk mengambil data dengan pagination
	dataQuery := baseQuery.Session(&gorm.Session{}).
		Preload("Journal").
		Preload("Journal.TransactionCategory"). // ✅ Tambahkan ini
		Order("journal_entries.date_inputed ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	if err := dataQuery.Find(&lines).Error; err != nil {
		return nil, 0, 0, err
	}

	var response []dto.JournalLineResponse
	for _, line := range lines {

		response = append(response, dto.JournalLineResponse{
			ID:                      line.JournalID.String(),
			Invoice:                 line.Journal.Invoice,
			Transaction_ID:          line.Journal.Transaction_ID,
			TransactionCategoryID:   line.Journal.TransactionCategoryID.String(),
			TransactionCategoryName: line.Journal.TransactionCategory.Name,
			DebitCategory:           line.Journal.TransactionCategory.DebitCategory,
			CreditCategory:          line.Journal.TransactionCategory.CreditCategory,
			Description:             line.Journal.Description,
			Partner:                 line.Journal.Partner,
			Amount:                  math.Abs(float64(line.Debit - line.Credit)),
			TransactionType:         string(line.TransactionType),
			DebitAccountType:        line.DebitAccountType,
			CreditAccountType:       line.CreditAccountType,
			Status:                  string(line.Journal.Status),
			CompanyID:               line.CompanyID.String(),
			DateInputed:             *line.Journal.DateInputed,
			DueDate:                 helper.SafeDueDate(line.Journal.DueDate),
			IsRepaid:                line.Journal.IsRepaid,
			Installment:             line.Journal.Installment,
			Note:                    line.Journal.Note,
		})
	}

	return response, totalExpense, total, nil
}
