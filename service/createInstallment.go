package service

import (
	"ayana/models"
	"fmt"
)

// func CreateInstallmentJournals(input models.JournalEntry) ([]models.JournalEntry, error) {
// 	var journals []models.JournalEntry

// 	// if err := c.ShouldBindJSON(&input); err != nil {
// 	// 	c.JSON(http.StatusBadRequest, gin.H{
// 	// 		"status":  "error",
// 	// 		"message": "Invalid input",
// 	// 		"details": err.Error(),
// 	// 	})
// 	// 	return
// 	// }

// 	// c.JSON(http.StatusOK, gin.H{
// 	// 	"status": "success",
// 	// 	"data":   input,
// 	// })

// 	// Hitung nilai cicilan per installment
// 	installmentAmount := input.Amount / int64(input.Installment)

// 	for i := 0; i < input.Installment; i++ {
// 		newJournal := models.JournalEntry{
// 			ID:                    uuid.New(),
// 			Invoice:               input.Invoice + "-" + string(rune(i+1)),
// 			Description:           input.Description,
// 			TransactionCategoryID: input.TransactionCategoryID,
// 			Amount:                installmentAmount,
// 			Partner:               input.Partner,
// 			TransactionType:       input.TransactionType,
// 			Status:                input.Status,
// 			CompanyID:             input.CompanyID,
// 			IsRepaid:              false,
// 			Installment:           input.Installment,
// 			Note:                  input.Note,
// 			DateInputed:           input.DateInputed,
// 			DueDate:               nil,
// 		}

// 		// DueDate per installment (misalnya tiap bulan)
// 		if input.DueDate != nil {
// 			due := input.DueDate.AddDate(0, i, 0) // add bulan
// 			newJournal.DueDate = &due
// 		}

// 		// Simpan ke DB
// 		if err := db.DB.Create(&newJournal).Error; err != nil {
// 			return nil, err
// 		}

// 		journals = append(journals, newJournal)
// 	}

// 	return journals, nil
// }

// func CreateInstallmentJournals(input models.JournalEntry) ([]models.JournalEntry, error) {

// 	// Debug 2: Retrieve transaction category details
// 	var trxCategory models.TransactionCategory
// 	// if err := db.DB.Preload("DebitAccount").Preload("CreditAccount").
// 	// 	First(&trxCategory, "id = ?", input.TransactionCategoryID).Error; err != nil {
// 	// 	return nil, fmt.Errorf("failed to retrieve transaction category: %v", err)
// 	// }

// 	// prettyJSON, err := json.MarshalIndent(trxCategory, "", "  ")
// 	// if err != nil {
// 	// 	fmt.Println("Error marshalling to JSON:", err)

// 	// }
// 	// fmt.Printf("Transaction Category:\n%s\n", prettyJSON)
// 	// Debug 3: Calculate amount per installment
// 	amountPerInstallment := input.Amount / int64(input.Installment)
// 	fmt.Printf("Amount Per Installment: %d\n", amountPerInstallment)

// 	// Debug 4: Prepare for installments and initialize journals array
// 	var journals []models.JournalEntry
// 	fmt.Println("Preparing to create installments...")

// 	// Debug 5: Loop through each installment
// 	for i := 0; i < input.Installment; i++ {
// 		// Debug 6: Log current installment index
// 		fmt.Printf("Creating Installment %d of %d\n", i+1, input.Installment)

// 		// Generate unique journal ID and calculate due date for installment
// 		journalID := uuid.New()
// 		dueDate := input.DateInputed.AddDate(0, i, 0)
// 		dueDatePtr := &dueDate

// 		// Debug 7: Log journal ID and due date for this installment
// 		fmt.Printf("Installment %d - Journal ID: %s, Due Date: %v\n", i+1, journalID, dueDate)

// 		// Determine debit amount (only the first installment gets the full debit amount)
// 		debitAmount := int64(0)
// 		if i == 0 {
// 			debitAmount = input.Amount
// 		}

// 		// Debug 8: Log debit amount for current installment
// 		fmt.Printf("Installment %d - Debit Amount: %d\n", i+1, debitAmount)

// 		// Create the journal entry for this installment
// 		journal := models.JournalEntry{
// 			ID:                    journalID,
// 			Invoice:               fmt.Sprintf("%s-%02d", input.Invoice, i+1),
// 			Description:           fmt.Sprintf("%s - Cicilan %d", input.Description, i+1),
// 			TransactionCategoryID: input.TransactionCategoryID,
// 			Amount:                amountPerInstallment,
// 			Partner:               input.Partner,
// 			TransactionType:       input.TransactionType,
// 			Status:                input.Status,
// 			DateInputed:           input.DateInputed,
// 			DueDate:               dueDatePtr,
// 			CompanyID:             input.CompanyID,
// 			CreatedAt:             time.Now(),
// 			UpdatedAt:             time.Now(),
// 		}

// 		// Debug 9: Log journal entry before adding lines
// 		fmt.Printf("Installment %d - Journal Entry: %+v\n", i+1, journal)

// 		// Create the journal lines for debit and credit
// 		journal.Lines = []models.JournalLine{
// 			{
// 				ID:          uuid.New(),
// 				JournalID:   journalID,
// 				AccountID:   trxCategory.DebitAccountID,
// 				Debit:       debitAmount,
// 				Credit:      0,
// 				Description: "Auto debit cicilan",
// 				CreatedAt:   time.Now(),
// 				UpdatedAt:   time.Now(),
// 			},
// 			{
// 				ID:          uuid.New(),
// 				JournalID:   journalID,
// 				AccountID:   trxCategory.CreditAccountID,
// 				Debit:       0,
// 				Credit:      amountPerInstallment,
// 				Description: "Auto credit cicilan",
// 				CreatedAt:   time.Now(),
// 				UpdatedAt:   time.Now(),
// 			},
// 		}

// 		// Debug 10: Log journal lines before adding to the journal
// 		fmt.Printf("Installment %d - Journal Lines: %+v\n", i+1, journal.Lines)

// 		// Add the created journal to the slice
// 		journals = append(journals, journal)
// 	}

// 	// Debug 11: Log the final journals slice before saving
// 	fmt.Printf("Final Journals: %+v\n", journals)

// 	// Commit the created journals to the database
// 	if err := db.DB.Create(&journals).Error; err != nil {
// 		return nil, fmt.Errorf("failed to create installment journals: %v", err)
// 	}

// 	// Debug 12: Confirm successful journal creation
// 	fmt.Println("Successfully created installment journals")

// 	// Return the created journals
// 	return journals, nil
// }

func CreateInstallmentJournals(input models.JournalEntry) ([]models.JournalEntry, error) {
	fmt.Println(">>> DEBUG input:", input) // Tambahkan log ini
	return []models.JournalEntry{input}, nil
}
