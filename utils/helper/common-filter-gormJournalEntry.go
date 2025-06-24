package helper

import (
	lib "ayana/lib"

	"gorm.io/gorm"
)

type JournalEntryFilterParams struct {
	DebitCategory  string
	CreditCategory string
	DateFilter     lib.DateFilter
	SortBy         string
	SortOrder      string
}

func ApplyCommonJournalEntryFiltersToGorm(query *gorm.DB, params JournalEntryFilterParams, applySort bool) (*gorm.DB, string, string) {
	// Filter kategori debit
	if params.DebitCategory != "" {
		query = query.Where("LOWER(transaction_categories.debit_category) = LOWER(?)", params.DebitCategory)
	}

	// Filter kategori kredit
	if params.CreditCategory != "" {
		query = query.Where("LOWER(transaction_categories.credit_category) = LOWER(?)", params.CreditCategory)
	}

	// Filter tanggal
	if params.DateFilter.StartDate != nil && params.DateFilter.EndDate != nil {
		query = query.Where("journal_entries.date_inputed BETWEEN ? AND ?", params.DateFilter.StartDate, params.DateFilter.EndDate)
	} else if params.DateFilter.StartDate != nil {
		query = query.Where("journal_entries.date_inputed >= ?", params.DateFilter.StartDate)
	} else if params.DateFilter.EndDate != nil {
		query = query.Where("journal_entries.date_inputed <= ?", params.DateFilter.EndDate)
	}

	// Sorting defaults
	sortBy := "date_inputed"
	if params.SortBy == "due_date" {
		sortBy = "due_date"
	}

	sortOrder := "asc"
	if params.SortOrder == "desc" {
		sortOrder = "desc"
	}

	// Apply ORDER BY hanya jika diminta
	if applySort {
		query = query.Order("journal_entries." + sortBy + " " + sortOrder)
	}

	return query, sortBy, sortOrder
}
