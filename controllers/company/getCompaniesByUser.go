package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetCompaniesByUser(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "parameter 'user_id' wajib diisi"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "user_id tidak valid", "error": err.Error()})
		return
	}

	var actor models.User
	if err := db.DB.First(&actor, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var companies []models.Company
	var total int64

	// build query: kalau bukan superadmin, filter via pivot user_companies
	query := db.DB.Model(&models.Company{})
	if actor.Role != "superadmin" {
		query = query.
			Joins("JOIN user_companies uc ON uc.company_id = companies.id").
			Where("uc.user_id = ?", actor.ID)
	}

	// hitung total
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghitung data", "error": err.Error()})
		return
	}

	// ambil data (preload relasi user, tanpa preload company)
	if err := query.
		Preload("Users.User").
		Order("CAST(companies.company_code AS INTEGER) ASC").
		Limit(limit).
		Offset(offset).
		Find(&companies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data perusahaan", "error": err.Error()})
		return
	}

	// DTO mapping
	type UserDTO struct {
		ID       uuid.UUID `json:"id"`
		UserID   uuid.UUID `json:"user_id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
	}

	type CompanyDTO struct {
		ID            uuid.UUID `json:"id"`
		Title         string    `json:"title"`
		CompanyCode   string    `json:"company_code"`
		Color         string    `json:"color"`
		HasProject    bool      `json:"has_project"`
		HasCustomer   bool      `json:"has_customer"`
		IsRetail      bool      `json:"is_retail"`
		CreatedAt     string    `json:"created_at"`
		UpdatedAt     string    `json:"updated_at"`
		AssignedUsers []UserDTO `json:"users"`
	}

	var companyList []CompanyDTO
	for _, comp := range companies {
		var userList []UserDTO
		for _, uc := range comp.Users {
			userList = append(userList, UserDTO{
				ID:       uc.ID,
				UserID:   uc.UserID,
				Username: uc.User.Username,
				Role:     uc.User.Role,
			})
		}

		companyList = append(companyList, CompanyDTO{
			ID:            comp.ID,
			Title:         comp.Title,
			CompanyCode:   comp.CompanyCode,
			Color:         comp.Color,
			HasProject:    comp.HasProject,
			HasCustomer:   comp.HasCustomer,
			IsRetail:      comp.IsRetail,
			CreatedAt:     comp.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     comp.UpdatedAt.Format(time.RFC3339),
			AssignedUsers: userList,
		})
	}

	// final response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"companyList":    companyList,
			"total_customer": total, // kalau ini maksudnya jumlah user per company, harus dihitung terpisah
			"page":           page,
			"limit":          limit,
			"total":          total,
		},
		"message": "Data perusahaan berhasil diambil",
		"status":  "sukses",
	})
}
