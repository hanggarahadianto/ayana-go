package routes

import (
	reservationController "ayana/controllers/reservation"

	"github.com/gin-gonic/gin"
)

func SetupReservationRouter(r *gin.Engine) {
	home := r.Group("/reservation")
	{
		home.GET("/get", reservationController.GetReservations)
		home.POST("/post/id", reservationController.CreateReservation)
	}
}
