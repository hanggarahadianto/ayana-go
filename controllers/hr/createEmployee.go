package controllers

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateEmployee(c *gin.Context) {
	var input models.Employee

	// Bind JSON ke struct input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi perusahaan
	if !helper.ValidateCompanyExist(input.CompanyID, c) {
		return
	}

	// Buat data karyawan baru
	employee := models.Employee{
		ID:                   uuid.New(), // Jika DB tidak handle auto UUID
		Name:                 input.Name,
		Address:              input.Address,
		Phone:                input.Phone,
		DateBirth:            input.DateBirth,
		MaritalStatus:        input.MaritalStatus,
		EmployeeEducation:    input.EmployeeEducation,
		Department:           input.Department,
		Gender:               input.Gender,
		Religion:             input.Religion, // Tambahkan field Religion
		Position:             input.Position,
		EmployeeStatus:       input.EmployeeStatus,
		EmployeeContractType: input.EmployeeContractType,
		CompanyID:            input.CompanyID,
	}

	// Simpan ke database
	if err := db.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   employee,
	})
}
