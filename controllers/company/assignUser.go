package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "parameter 'user_id' wajib"})
		return
	}
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "user_id tidak valid"})
		return
	}

	var actor models.User
	if err := db.DB.First(&actor, "id = ?", actorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Actor tidak ditemukan"})
		return
	}

	if actor.Role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Hanya superadmin"})
		return
	}

	var req AssignCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Request tidak valid", "error": err.Error()})
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "company_id tidak valid"})
		return
	}

	var company models.Company
	if err := db.DB.First(&company, "id = ?", companyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Company tidak ditemukan"})
		return
	}

	// Ambil relasi lama
	var existing []models.UserCompany
	db.DB.Where("company_id = ?", companyID).Find(&existing)

	// Buat map untuk cek cepat
	keepMap := make(map[uuid.UUID]bool)
	for _, idStr := range req.UserIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			keepMap[id] = true
		}
	}

	// Hapus relasi lama yang tidak ada di req.UserIDs
	for _, e := range existing {
		if !keepMap[e.UserID] {
			db.DB.Delete(&e)
		}
	}

	// Tambah relasi baru
	created := []string{}
	skipped := []string{}
	for _, uidStr := range req.UserIDs {
		uID, err := uuid.Parse(uidStr)
		if err != nil {
			skipped = append(skipped, uidStr)
			continue
		}

		// Cek user exist
		var u models.User
		if err := db.DB.First(&u, "id = ?", uID).Error; err != nil {
			skipped = append(skipped, uidStr)
			continue
		}

		// Cek apakah sudah ada
		var count int64
		db.DB.Model(&models.UserCompany{}).Where("user_id = ? AND company_id = ?", uID, companyID).Count(&count)
		if count > 0 {
			continue
		}

		uc := models.UserCompany{
			UserID:    uID,
			CompanyID: companyID,
		}
		if err := db.DB.Create(&uc).Error; err == nil {
			created = append(created, uidStr)
		} else {
			skipped = append(skipped, uidStr)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Assign selesai",
		"created": created,
		"skipped": skipped,
	})
}
