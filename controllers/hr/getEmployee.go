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
	var employees []models.Employee
	var total int64
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
	isAgent := c.Query("is_agent")

	// 🔁 Build query awal
	query := db.DB.Model(&models.Employee{}).Where("company_id = ?", companyID)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// ✅ Apply filter is_agent SEBELUM count
	switch isAgent {
	case "true":
		query = query.Where("is_agent = ?", true)
	case "false":
		query = query.Where("is_agent = ?", false)
	}

	// 🔢 Hitung total setelah semua filter
	if err := query.Count(&total).Error; err != nil {
		log.Println("🔴 Count error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data karyawan"})
		return
	}

	// 📦 Ambil data
	if err := query.
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Order("name ASC"). // urut berdasarkan nama
		Find(&employees).Error; err != nil {
		log.Println("🔴 Find error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data karyawan"})
		return
	}

	// ✅ Response
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
