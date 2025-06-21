package controllers

import (
	"ayana/models"
	"ayana/service"
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

	updatedCustomer, err := service.UpdateCustomerService(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCustomer)
}
