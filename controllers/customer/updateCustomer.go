package controllers

import (
	"ayana/models"
	customer "ayana/service/customer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var input models.Customer

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCustomer, err := customer.UpdateCustomerService(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCustomer)
}
