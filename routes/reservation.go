package routes

import (
	reservationController "ayana/controllers/reservation"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupReservationRouter(r *gin.Engine) {
	reservation := r.Group("/reservation", middlewares.AuthMiddleware())
	{
		reservation.GET("/get", reservationController.GetReservations)
		reservation.POST("/post", reservationController.CreateReservation) // Dynamic :id
	}
}
