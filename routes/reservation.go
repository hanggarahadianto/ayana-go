package routes

import (
	reservationController "ayana/controllers/reservation"

	"github.com/gin-gonic/gin"
)

func SetupReservationRouter(r *gin.Engine) {
	reservation := r.Group("/reservation")
	{
		reservation.GET("/get", reservationController.GetReservations)
		reservation.POST("/post/:id", reservationController.CreateReservation) // Dynamic :id
	}
}
