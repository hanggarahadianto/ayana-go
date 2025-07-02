package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeleteEmployee menghapus data karyawan berdasarkan ID
func DeleteEmployee(c *gin.Context) {
	employeeIDParam := c.Param("id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID karyawan tidak valid"})
		return
	}

	// Periksa apakah karyawan ada
	var employee models.Employee
	if err := db.DB.First(&employee, "id = ?", employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data karyawan tidak ditemukan"})
		return
	}

	// Hapus data karyawan
	if err := db.DB.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data karyawan berhasil dihapus",
	})
}
