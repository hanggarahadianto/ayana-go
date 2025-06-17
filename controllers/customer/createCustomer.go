package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCustomer(c *gin.Context) {
	var input models.Customer

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := models.Customer{
		ID:            uuid.New(),
		Name:          input.Name,
		Address:       input.Address,
		Phone:         input.Phone, // Tambahkan ini
		Status:        input.Status,
		PaymentMethod: input.PaymentMethod,
		DateInputed:   input.DateInputed,
		Amount:        input.Amount,
		Marketer:      input.Marketer,
		HomeID:        input.HomeID,
		ProductUnit:   input.ProductUnit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}
