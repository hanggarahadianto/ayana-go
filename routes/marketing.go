package routes

import (
	marketingController "ayana/controllers/marketing"

	"github.com/gin-gonic/gin"
)

func SetupMarketingRouter(r *gin.Engine) {
	marketing := r.Group("/marketing")
	{
		marketing.GET("/get", marketingController.GetMarketerPerformanceHandler)
		// marketing.POST("/post", marketingController.CreateReservation) // Dynamic :id
	}
}
