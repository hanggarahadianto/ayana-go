package routes

import (
	marketingController "ayana/controllers/marketing"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupMarketingRouter(r *gin.Engine) {
	marketing := r.Group("/marketing", middlewares.AuthMiddleware())
	{
		marketing.GET("/get", marketingController.GetMarketerPerformanceHandler)
		// marketing.POST("/post", marketingController.CreateReservation) // Dynamic :id
	}
}
