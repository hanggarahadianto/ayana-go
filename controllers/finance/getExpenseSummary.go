package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetExpenseSummary(c *gin.Context) {
	companyID := c.Query("company_id")

	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	// Ambil page & limit dari query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if err := db.DB.Model(&models.JournalLine{}).
		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("accounts.company_id = ? AND accounts.type = ? AND journal_entries.status = ?", companyID, "Expense (Beban)", "paid").
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count expenses"})
		return
	}

	// Ambil data dengan limit dan offset
	var journalLines []models.JournalLine
	if err := db.DB.Preload("Account").
		Joins("JOIN accounts ON accounts.id = journal_lines.account_id").
		Joins("JOIN journal_entries ON journal_entries.id = journal_lines.journal_id").
		Where("accounts.company_id = ? AND accounts.type = ? AND journal_entries.status = ?", companyID, "Expense (Beban)", "paid").
		Limit(limit).Offset(offset).
		Find(&journalLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expense data"})
		return
	}
	var totalExpense int64 = 0
	var expenseList []map[string]interface{}

	for _, line := range journalLines {
		balance := line.Debit - line.Credit
		totalExpense += balance

		expenseList = append(expenseList, map[string]interface{}{
			"id":          line.ID,
			"account":     line.Account.Name,
			"category":    line.Account.Category,
			"description": line.Description,
			"amount":      balance,
			"date":        line.CreatedAt.Format("2006-01-02"),
			// "status":      "paid", // opsional
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"expenseList":   expenseList,
			"total_expense": totalExpense,
			"page":          page,
			"limit":         limit,
			"total":         total,
		},
		"message": "Expense summary retrieved successfully",
		"status":  "success",
	})
}
