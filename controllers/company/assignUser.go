package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignCompanyRequest struct {
	CompanyID string   `json:"company_id" binding:"required"`
	UserIDs   []string `json:"user_ids" binding:"required"`
}

// GET /company/get?user_id=<uuid>&page=1&limit=10

// POST /company/post/assign-user?user_id=<actor_uuid>
func AssignCompanyToUsers(c *gin.Context) {
	actorIDStr := c.Query("user_id")
	if actorIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "parameter 'user_id' (yang melakukan aksi) wajib diisi"})
		return
	}
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "user_id tidak valid", "error": err.Error()})
		return
	}

	var actor models.User
	if err := db.DB.First(&actor, "id = ?", actorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User (actor) tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// hanya superadmin boleh assign
	if actor.Role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Hanya superadmin yang bisa assign company"})
		return
	}

	// bind request
	var req AssignCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Request body tidak valid", "error": err.Error()})
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "company_id tidak valid", "error": err.Error()})
		return
	}

	// cek company ada
	var company models.Company
	if err := db.DB.First(&company, "id = ?", companyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Company tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	created := []string{}
	skipped := []string{}
	for _, uidStr := range req.UserIDs {
		uID, err := uuid.Parse(uidStr)
		if err != nil {
			skipped = append(skipped, uidStr) // invalid uuid
			continue
		}

		// optional: cek user exist
		var u models.User
		if err := db.DB.First(&u, "id = ?", uID).Error; err != nil {
			skipped = append(skipped, uidStr)
			continue
		}

		var count int64
		db.DB.Model(&models.UserCompany{}).
			Where("user_id = ? AND company_id = ?", uID, companyID).
			Count(&count)
		if count > 0 {
			skipped = append(skipped, uidStr) // sudah ada relasi
			continue
		}

		uc := models.UserCompany{
			UserID:    uID,
			CompanyID: companyID,
		}
		if err := db.DB.Create(&uc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal insert relasi", "error": err.Error()})
			return
		}
		created = append(created, uidStr)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Assign selesai",
		"created": created,
		"skipped": skipped,
	})
}
