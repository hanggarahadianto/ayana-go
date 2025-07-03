package controllers

import (
	"ayana/db"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateEmployee(c *gin.Context) {
	employeeIDParam := c.Param("id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah karyawan ada
	var existingEmployee models.Employee
	if err := db.DB.First(&existingEmployee, "id = ?", employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Validasi Company
	if !helper.ValidateCompanyExist(input.CompanyID, c) {
		return
	}

	// Update field
	existingEmployee.Name = input.Name
	existingEmployee.Address = input.Address
	existingEmployee.Phone = input.Phone
	existingEmployee.DateBirth = input.DateBirth
	existingEmployee.MaritalStatus = input.MaritalStatus
	existingEmployee.EmployeeEducation = input.EmployeeEducation
	existingEmployee.Department = input.Department
	existingEmployee.Gender = input.Gender
	existingEmployee.Religion = input.Religion
	existingEmployee.Position = input.Position
	existingEmployee.EmployeeStatus = input.EmployeeStatus
	existingEmployee.EmployeeContractType = input.EmployeeContractType
	existingEmployee.CompanyID = input.CompanyID

	if err := db.DB.Save(&existingEmployee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   existingEmployee,
	})
}
