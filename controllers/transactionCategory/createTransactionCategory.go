package controller

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransactionCategory(c *gin.Context) {
	var input models.TransactionCategory

	// Bind the JSON input into the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Validate input data
	if input.Name == "" || input.Category == "" || input.Description == "" || input.CompanyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Missing required fields",
		})
		return
	}

	// Check if debit account exists
	var debitAccount models.Account
	if err := db.DB.First(&debitAccount, "id = ?", input.DebitAccountID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Debit account not found",
		})
		return
	}

	// Check if credit account exists
	var creditAccount models.Account
	if err := db.DB.First(&creditAccount, "id = ?", input.CreditAccountID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Credit account not found",
		})
		return
	}

	// Create the transaction category
	transactionCategory := models.TransactionCategory{
		ID:                uuid.New(),
		Name:              input.Name,
		DebitAccountID:    input.DebitAccountID,
		DebitAccountType:  input.DebitAccountType,
		CreditAccountID:   input.CreditAccountID,
		CreditAccountType: input.CreditAccountType,
		Category:          input.Category,
		Description:       input.Description,
		CompanyID:         input.CompanyID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save the transaction category to the database
	if err := db.DB.Create(&transactionCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create transaction category",
			"details": err.Error(),
		})
		return
	}

	if err := db.DB.Preload("DebitAccount").Preload("CreditAccount").First(&transactionCategory, "id = ?", transactionCategory.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to load related accounts",
		})
		return
	}

	// Respond with the newly created transaction category
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactionCategory,
	})
}
