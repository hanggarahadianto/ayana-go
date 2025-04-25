package helper

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ValidateCompanyID memvalidasi apakah CompanyID ada di database
func ValidateCompanyID(companyID uuid.UUID, c *gin.Context) bool {
	var company models.Company
	if err := db.DB.First(&company, "id = ?", companyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid CompanyID"})
		return false
	}
	return true
}
