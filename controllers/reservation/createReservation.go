package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateReservation(c *gin.Context) {
	id := c.Param("id")
	var reservationData models.Reservation

	if err := c.ShouldBindJSON(&reservationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	newReservation := models.Reservation{
		Name:      reservationData.Name,
		Phone:     reservationData.Phone,
		Date:      reservationData.Date,
		Time:      reservationData.Time,
		Home_ID:   id,
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
