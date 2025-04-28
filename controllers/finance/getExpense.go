package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetExpenseSummary mengembalikan ringkasan pengeluaran atau daftar pengeluaran
func GetExpenseSummary(c *gin.Context) {

	// Validasi parameter wajib
	companyIDStr := c.Query("company_id")

	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	// Ambil filter tanggal
	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	// Cek apakah hanya ingin summary
	if c.Query("summary_only") == "true" {
		totalExpense, err := service.GetExpenseSummaryOnly(companyIDStr, dateFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total pengeluaran"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"total_expense": totalExpense,
			},
			"message": "Ringkasan pengeluaran berhasil diambil",
			"status":  "sukses",
		})
		return
	}

	// Jika ingin daftar pengeluaran beserta summary
	pagination := helper.GetPagination(c)

	params := service.ExpenseFilterParams{
		CompanyID:  companyID.String(),
		Pagination: pagination,
		DateFilter: dateFilter,
	}

	data, totalExpense, total, err := service.GetExpenses(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengeluaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"expenseList":   data,
			"total_expense": totalExpense,
			"page":          pagination.Page,
			"limit":         pagination.Limit,
			"total":         total,
		},
		"message": "Pengeluaran berhasil diambil",
		"status":  "sukses",
	})
}
