package helper

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ValidateCompanyID memvalidasi apakah CompanyID ada di database
func ValidateCompanyExist(companyID uuid.UUID, c *gin.Context) bool {
	var company models.Company
	if err := db.DB.First(&company, "id = ?", companyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Company not found"})
		return false
	}
	return true
}

// ValidateAndParseCompanyID memvalidasi dan memparsing company_id menjadi UUID
func ValidateAndParseCompanyID(companyIDStr string, c *gin.Context) (uuid.UUID, bool) {
	// Pastikan company_id tidak kosong
	if companyIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return uuid.Nil, false
	}

	// Parse ke UUID
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for Company ID"})
		return uuid.Nil, false
	}

	// Validasi company_id di database
	var company models.Company
	if err := db.DB.First(&company, "id = ?", companyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Company not found"})
		return uuid.Nil, false
	}

	return companyID, true
}
