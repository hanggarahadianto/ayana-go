package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateCompany(c *gin.Context) {

	var companyData models.Company

	if err := c.ShouldBindJSON(&companyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	newCompany := models.Company{
		Title:       companyData.Title,
		CompanyCode: companyData.CompanyCode,

		CreatedAt: now,
		UpdatedAt: now,
	}

	result := db.DB.Create(&newCompany)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newCompany,
	})

}
