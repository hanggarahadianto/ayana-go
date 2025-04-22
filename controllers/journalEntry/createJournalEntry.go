package controller

import (
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"

	Service "ayana/service"
)

func CreateJournalEntry(c *gin.Context) {
	var input models.JournalEntry

	// prettyJSON, err := json.MarshalIndent(input, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error marshalling to JSON:", err)

	// }

	// fmt.Println(prettyJSON)

	// if input.TransactionCategoryID == uuid.Nil || input.Amount <= 0 || input.CompanyID == uuid.Nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Missing required fields",
	// 	})
	// 	return
	// }

	// var company models.Company
	// if err := db.DB.First(&company, "id = ?", input.CompanyID).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Company doest exist",
	// 	})
	// 	return
	// }

	// // Validate Transaction Category
	// var trxCategory models.TransactionCategory
	// if err := db.DB.Preload("DebitAccount").Preload("CreditAccount").First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Transaction category not found",
	// 	})
	// 	return
	// }

	// // Validate Accounts associated with the Transaction Category
	// var debitAccount models.Account
	// if err := db.DB.First(&debitAccount, "id = ?", trxCategory.DebitAccountID).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Debit account not found",
	// 	})
	// 	return
	// }

	// var creditAccount models.Account
	// if err := db.DB.First(&creditAccount, "id = ?", trxCategory.CreditAccountID).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Credit account not found",
	// 	})
	// 	return
	// }

	if input.Installment > 1 {
		journals, err := Service.CreateInstallmentJournals(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to create installment journals",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   journals,
		})
		return
	}

	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Invalid input",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "success",
	// 	"data":   input,
	// })

	// Create JournalEntry
	// journal := models.JournalEntry{
	// 	ID:                    uuid.New(),
	// 	Invoice:               input.Invoice,
	// 	Description:           input.Description,
	// 	TransactionCategoryID: input.TransactionCategoryID,
	// 	Amount:                input.Amount,
	// 	Partner:               input.Partner,
	// 	TransactionType:       input.TransactionType,
	// 	Status:                input.Status,
	// 	DateInputed:           input.DateInputed,
	// 	DueDate:               input.DueDate,
	// 	Note:                  input.Note,
	// 	CompanyID:             input.CompanyID,
	// 	CreatedAt:             time.Now(),
	// 	UpdatedAt:             time.Now(),
	// }

	// // Create Lines (Auto from Transaction Category)
	// journal.Lines = []models.JournalLine{
	// 	{
	// 		ID:          uuid.New(),
	// 		JournalID:   journal.ID,
	// 		AccountID:   trxCategory.DebitAccountID,
	// 		CompanyID:   journal.CompanyID,
	// 		Debit:       input.Amount,
	// 		Credit:      0,
	// 		Description: input.Description,

	// 		CreatedAt: time.Now(),
	// 		UpdatedAt: time.Now(),
	// 	},
	// 	{
	// 		ID:          uuid.New(),
	// 		JournalID:   journal.ID,
	// 		AccountID:   trxCategory.CreditAccountID,
	// 		CompanyID:   journal.CompanyID,
	// 		Debit:       0,
	// 		Credit:      input.Amount,
	// 		Description: input.Description,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	},
	// }

	// if err := db.DB.Create(&journal).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  "error",
	// 		"message": "Failed to create journal entry",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }

	// var journalWithLinesAndCategory models.JournalEntry
	// if err := db.DB.Preload("Lines.Account").
	// 	Preload("TransactionCategory").
	// 	Preload("TransactionCategory.DebitAccount").
	// 	Preload("TransactionCategory.CreditAccount").
	// 	First(&journalWithLinesAndCategory, "id = ?", journal.ID).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  "error",
	// 		"message": "Failed to preload journal lines and transaction category",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "success",
	// 	"data":   journalWithLinesAndCategory,
	// })

}
