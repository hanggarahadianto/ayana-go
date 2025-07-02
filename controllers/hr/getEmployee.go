package controllers

import (
	"ayana/db"
	lib "ayana/lib"

	"ayana/models"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
	pagination := lib.GetPagination(c)
	if !lib.ValidatePagination(pagination, c) {
		return
	}

	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company Id wajib diisi"})
		return
	}

	search := c.Query("search")
	var employees []models.Employee
	var total int64

	query := db.DB.Model(&models.Employee{}).Where("company_id = ?", companyID)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		log.Println("🔴 Count error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data karyawan"})
		return
	}

	if err := query.
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Order("COALESCE(updated_at, created_at) DESC").
		Find(&employees).Error; err != nil {
		log.Println("🔴 Find error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"employeeList":   employees,
			"total_employee": total,
			"page":           pagination.Page,
			"limit":          pagination.Limit,
			"total":          total,
		},
		"message": "Data karyawan berhasil diambil",
		"status":  "sukses",
	})
}
