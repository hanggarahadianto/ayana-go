// package controller

// import (
// 	"ayana/db"
// 	"ayana/models"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// func GetCashSummary(c *gin.Context) {
// 	companyID := c.Query("company_id")
// 	if companyID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
// 		return
// 	}

// 	transactionTypeFilter := strings.ToLower(c.Query("transaction_type"))

// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 10
// 	}
// 	offset := (page - 1) * limit

// 	// ➕ Tambahan: Query full journalLines untuk total summary
// 	var allJournalLines []models.JournalLine
// 	if err := db.DB.
// 		Preload("Account").
// 		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
// 		Where("accounts.company_id = ?", companyID).
// 		Find(&allJournalLines).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate summary data"})
// 		return
// 	}

// 	var journalLines []models.JournalLine
// 	if err := db.DB.
// 		Preload("Account").
// 		Preload("Journal").
// 		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
// 		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
// 		Where("accounts.company_id = ?", companyID).
// 		Where("LOWER(journal_entries.transaction_type) = ?", transactionTypeFilter).
// 		Limit(limit).
// 		Offset(offset).
// 		Find(&journalLines).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve journal data"})
// 		return
// 	}

// 	var cashList []map[string]interface{}
// 	for _, line := range journalLines {
// 		account := line.Account
// 		balance := line.Debit - line.Credit
// 		transactionType := line.TransactionType

// 		if string(transactionType) == "" {
// 			if balance > 0 {
// 				transactionType = models.TransactionType("payin")
// 			} else if balance < 0 {
// 				transactionType = models.TransactionType("payout")
// 			}
// 		}

// 		transactionTypeStr := strings.ToLower(string(transactionType))
// 		if transactionTypeFilter != "" && transactionTypeStr != transactionTypeFilter {
// 			continue
// 		}

// 		if line.Journal.Status != "unpaid" {
// 			continue
// 		}

// 		cashList = append(cashList, map[string]interface{}{
// 			"id":               line.ID,
// 			"description":      account.Name,
// 			"amount":           balance,
// 			"date":             line.CreatedAt.Format("2006-01-02"),
// 			"status":           line.Journal.Status,
// 			"transaction_type": transactionType,
// 			"note":             line.Journal.Note,
// 		})
// 	}

// 	// total := int64(len(cashList))

// 	total := int64(len(cashList))

// 	// Hitung total_cashin dari cashList
// 	totalCashIn := int64(0)
// 	for _, item := range cashList {
// 		if amount, ok := item["amount"].(int64); ok {
// 			totalCashIn += amount
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": gin.H{
// 			"cashList":     cashList,
// 			"total_cashin": totalCashIn,

// 			"page":  page,
// 			"limit": limit,
// 			"total": total,
// 		},
// 		"message": "Cash summary retrieved successfully",
// 		"status":  "success",
// 	})
// }

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

	transactionTypeFilter := strings.ToLower(c.Query("transaction_type"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query all journalLines for totalCashIn (unfiltered, unpaid, payin only)
	var allJournalLines []models.JournalLine
	if err := db.DB.
		Preload("Account").
		Preload("Journal").
		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("accounts.company_id = ?", companyID).
		Where("journal_entries.status = ?", "unpaid").
		Where("LOWER(journal_entries.transaction_type) = ?", "payin").
		Find(&allJournalLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate summary data"})
		return
	}

	// Calculate totalCashIn from all journal lines (only payin)
	totalCashIn := int64(0)
	for _, line := range allJournalLines {
		balance := line.Debit - line.Credit
		if balance > 0 { // Only cash in (payin)
			totalCashIn += balance
		}
	}

	// Query filtered journalLines for pagination and transactionType
	var journalLines []models.JournalLine
	if err := db.DB.
		Preload("Account").
		Preload("Journal").
		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("accounts.company_id = ?", companyID).
		Where("LOWER(journal_entries.transaction_type) = ?", transactionTypeFilter).
		Limit(limit).
		Offset(offset).
		Find(&journalLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve journal data"})
		return
	}

	var cashList []map[string]interface{}
	for _, line := range journalLines {
		account := line.Account
		balance := line.Debit - line.Credit
		transactionType := line.TransactionType

		if string(transactionType) == "" {
			if balance > 0 {
				transactionType = models.TransactionType("payin")
			} else if balance < 0 {
				transactionType = models.TransactionType("payout")
			}
		}

		transactionTypeStr := strings.ToLower(string(transactionType))
		if transactionTypeFilter != "" && transactionTypeStr != transactionTypeFilter {
			continue
		}

		if line.Journal.Status != "unpaid" {
			continue
		}

		cashList = append(cashList, map[string]interface{}{
			"id":               line.ID,
			"description":      account.Name,
			"amount":           balance,
			"date":             line.CreatedAt.Format("2006-01-02"),
			"status":           line.Journal.Status,
			"transaction_type": transactionType,
			"note":             line.Journal.Note,
		})
	}

	// Calculate total of filtered cashList (if needed for pagination)
	total := int64(len(cashList))

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"cashList":     cashList,
			"total_cashin": totalCashIn, // Total Cash In from all journal lines (payin only)
			"page":         page,
			"limit":        limit,
			"total":        total,
		},
		"message": "Cash summary retrieved successfully",
		"status":  "success",
	})
}
