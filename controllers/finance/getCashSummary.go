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
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	transactionTypeFilter := c.Query("transaction_type")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var journalLines []models.JournalLine
	if err := db.DB.Preload("Account").Preload("Journal").
		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
		Where("accounts.company_id = ?", companyID).
		Limit(limit).Offset(offset).
		Find(&journalLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve journal data"})
		return
	}

	var cashList []map[string]interface{}
	netAssets := int64(0)
	availableCash := int64(0)

	for _, line := range journalLines {
		account := line.Account
		balance := line.Debit - line.Credit
		transactionType := line.TransactionType

		category := ""
		if balance > 0 {
			category = "inflow" // atau "penerimaan"
		} else if balance < 0 {
			category = "outflow" // atau "pengeluaran"
		}

		if account.Type == "Asset (Aset)" {
			netAssets += balance
		} else if account.Type == "Liability (Kewajiban)" {
			netAssets -= balance
		}

		if account.Type == "Asset (Aset)" && (account.Name == "Kas" || account.Name == "Bank") {
			availableCash += balance
		}

		if transactionTypeFilter != "" && !strings.EqualFold(string(transactionType), transactionTypeFilter) {
			continue
		}

		cashList = append(cashList, map[string]interface{}{
			"id":               line.ID,
			"description":      account.Name,
			"amount":           balance,
			"date":             line.CreatedAt.Format("2006-01-02"),
			"status":           "unpaid",
			"transaction_type": transactionType, // payin / payout
			"note":             line.Journal.Note,
			"category":         category,
		})
	}
	total := int64(len(cashList))

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"cashList":       cashList,
			"available_cash": availableCash,
			"net_assets":     netAssets,
			"page":           page,
			"limit":          limit,
			"total":          total, // ✅ Sudah sesuai hasil filter
		},
		"message": "Cash summary retrieved successfully",
		"status":  "success",
	})

}
