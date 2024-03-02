package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateMarketing(c *gin.Context) {

	var marketingData models.Marketing

	if err := c.ShouldBindJSON(&marketingData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	newMarketing := models.Marketing{
		Name:  marketingData.Name,
		Phone: marketingData.Phone,

		CreatedAt: now,
		UpdatedAt: now,
	}

	result := db.DB.Debug().Create(&newMarketing)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newMarketing,
	})

}
