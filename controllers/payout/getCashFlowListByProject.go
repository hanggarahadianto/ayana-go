package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCashFlowListByProjectId(c *gin.Context) {

	projectId := c.Param("id")

	// Define a variable to hold the weekly progress with related materials and workers
	var cashFlowList []models.CashFlow

	// Query the database with Preload to include materials and workers

	result := db.DB.Debug().
		Preload("Good").
		Where("project_id = ?", projectId).
		Find(&cashFlowList) // Use Find instead of First

	// Check if the query resulted in an error
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Cash Flow not found",
		})
		return
	}

	// Return the weekly progress with materials and workers
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   cashFlowList,
	})
}
