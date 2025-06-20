package controllers

import (
	"ayana/models"
	"ayana/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCustomer(c *gin.Context) {

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := service.CreateCustomer(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}
