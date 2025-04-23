package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCashSummary(c *gin.Context) {
	// Mendapatkan company_id dari query parameter
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	// Mendapatkan transaction_type dari query parameter
	transactionTypeFilter := strings.ToLower(c.Query("transaction_type"))

	// Mendapatkan page dan limit dari query parameter dengan nilai default
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query ke database untuk mendapatkan journal lines
	var journalLines []models.JournalLine
	if err := db.DB.Preload("Account").Preload("Journal").
		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
		Where("accounts.company_id = ?", companyID).
		Limit(limit).Offset(offset).
		Find(&journalLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve journal data"})
		return
	}

	// Menyimpan data cash summary
	var cashList []map[string]interface{}
	netAssets := int64(0)
	availableCash := int64(0)

	// Looping melalui setiap journal line untuk kalkulasi dan filter
	for _, line := range journalLines {
		account := line.Account
		balance := line.Debit - line.Credit
		transactionType := line.TransactionType

		// Menambahkan ke netAssets berdasarkan jenis akun
		if account.Type == "Asset (Aset)" {
			netAssets += balance
		} else if account.Type == "Liability (Kewajiban)" {
			netAssets -= balance
		}

		// Menambahkan ke availableCash jika akun adalah Kas atau Bank
		if account.Type == "Asset (Aset)" && (account.Name == "Kas" || account.Name == "Bank") {
			availableCash += balance
		}

		// Mengatur default transactionType jika kosong
		if string(transactionType) == "" {
			if balance > 0 {
				transactionType = models.TransactionType("payin")
			} else if balance < 0 {
				transactionType = models.TransactionType("payout")
			}
		}

		// Filter berdasarkan transactionType jika ada filter
		transactionTypeStr := strings.ToLower(string(transactionType))
		if transactionTypeFilter != "" && transactionTypeStr != transactionTypeFilter {
			continue
		}

		// Menambahkan data ke cashList
		cashList = append(cashList, map[string]interface{}{
			"id":               line.ID,
			"description":      account.Name,
			"amount":           balance,
			"date":             line.CreatedAt.Format("2006-01-02"),
			"status":           "unpaid",
			"transaction_type": transactionType, // payin / payout
			"note":             line.Journal.Note,
		})
	}

	// Menentukan total data yang didapat
	total := int64(len(cashList))

	// Mengirimkan response JSON
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"cashList":       cashList,
			"available_cash": availableCash,
			"net_assets":     netAssets,
			"page":           page,
			"limit":          limit,
			"total":          total,
		},
		"message": "Cash summary retrieved successfully",
		"status":  "success",
	})
}
