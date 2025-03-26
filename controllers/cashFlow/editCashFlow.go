package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateCashFlow(c *gin.Context) {
	var cashFlow models.CashFlow

	// Ambil ID dari parameter URL
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Cari CashFlow berdasarkan ID
	if err := db.DB.First(&cashFlow, "id = ?", parsedID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "CashFlow not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve CashFlow"})
		}
		return
	}

	// Bind data yang dikirim oleh client ke struct CashFlow
	if err := c.ShouldBindJSON(&cashFlow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi apakah ProjectID yang baru valid
	var project models.Project
	if err := db.DB.Where("id = ?", cashFlow.ProjectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking Project ID"})
		}
		return
	}

	// Update data
	cashFlow.UpdatedAt = time.Now()
	if err := db.DB.Save(&cashFlow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update CashFlow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   cashFlow,
	})
}
