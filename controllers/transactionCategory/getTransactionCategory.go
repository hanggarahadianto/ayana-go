package controller

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTransactionCategory(c *gin.Context) {
	id := c.Query("id")
	companyID := c.Query("company_id")
	transactionType := c.Query("transaction_type") // payin / payout

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var transactions []models.TransactionCategory
	var total int64

	tx := db.DB.Model(&models.TransactionCategory{}).
		Joins("JOIN accounts AS debit ON debit.id = transaction_categories.debit_account_id").
		Where("transaction_categories.company_id = ?", companyID)

	if id != "" {
		tx = tx.Where("transaction_categories.id = ?", id)
	}

	if transactionType == "payin" {
		tx = tx.Where("LOWER(debit.type) LIKE ?", "asset%")
	} else if transactionType == "payout" {
		tx = tx.Joins("JOIN accounts AS credit ON credit.id = transaction_categories.credit_account_id").
			Where("LOWER(credit.type) LIKE ?", "asset%")
	}

	if err := tx.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transaction categories"})
		return
	}

	if err := tx.Preload("DebitAccount").
		Preload("CreditAccount").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction categories"})
		return
	}

	// 🔁 Mapping ke response
	var responses []dto.TransactionCategoryResponse
	for _, t := range transactions {
		res := dto.TransactionCategoryResponse{
			ID:                t.ID,
			Name:              t.Name,
			DebitAccountID:    t.DebitAccountID,
			DebitAccountType:  t.DebitAccountType,
			CreditAccountID:   t.CreditAccountID,
			CreditAccountType: t.CreditAccountType,
			Category:          t.Category,
			Description:       t.Description,
			CompanyID:         t.CompanyID,
			DebitAccount: dto.AccountResponse{
				ID:          t.DebitAccount.ID,
				Code:        t.DebitAccount.Code,
				Name:        t.DebitAccount.Name,
				Type:        t.DebitAccount.Type,
				Category:    t.DebitAccount.Category,
				Description: t.DebitAccount.Description,
				CompanyID:   t.DebitAccount.CompanyID,
			},
			CreditAccount: dto.AccountResponse{
				ID:          t.CreditAccount.ID,
				Code:        t.CreditAccount.Code,
				Name:        t.CreditAccount.Name,
				Type:        t.CreditAccount.Type,
				Category:    t.CreditAccount.Category,
				Description: t.CreditAccount.Description,
				CompanyID:   t.CreditAccount.CompanyID,
			},
		}
		responses = append(responses, res)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   responses,
	})
}
