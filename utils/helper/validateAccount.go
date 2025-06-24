package helper

import (
	lib "ayana/lib"
	"ayana/models"
	"fmt"
)

// ValidateAccount untuk memvalidasi tipe dan kategori akun
func ValidateAccount(account *models.Account) error {
	// Validasi type akun
	validTypes := map[string]bool{
		"Asset":     true,
		"Liability": true,
		"Equity":    true,
		"Revenue":   true,
		"Expense":   true,
	}

	if _, valid := validTypes[account.Type]; !valid {
		return fmt.Errorf("invalid account type: %s", account.Type)
	}

	// Validasi kategori berdasarkan type akun
	if categories, ok := lib.ValidCategories[account.Type]; ok {
		valid := false
		for _, category := range categories {
			if account.Category == category {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid category for %s type: %s", account.Type, account.Category)
		}
	} else {
		return fmt.Errorf("invalid category type: %s", account.Type)
	}

	return nil
}
