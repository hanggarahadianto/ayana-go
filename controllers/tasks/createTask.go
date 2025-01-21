package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {

	var reservationData models.Reservation

	if err := c.ShouldBindJSON(&reservationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	newReservation := models.Reservation{
		Name:      reservationData.Name,
		Email:     reservationData.Email,
		Phone:     reservationData.Phone,
		HomeID:    reservationData.HomeID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := db.DB.Debug().Create(&newReservation)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newReservation,
	})

}
