package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReservations(c *gin.Context) {
	var reservationList []models.Reservation

	result := db.DB.Debug().Order("created_at desc, updated_at desc").Find(&reservationList)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reservationList,
	})

}
