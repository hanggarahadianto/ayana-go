package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var input models.Customer

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer models.Customer

	// Check if customer exists
	if err := db.DB.First(&customer, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Update fields
	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Status = input.Status
	customer.PaymentMethod = input.PaymentMethod
	customer.Amount = input.Amount
	customer.DateInputed = input.DateInputed
	customer.Marketer = input.Marketer
	customer.HomeID = input.HomeID
	customer.UpdatedAt = time.Now()

	if err := db.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
