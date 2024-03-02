package routes

import (
	marketingController "ayana/controllers/marketing"

	"github.com/gin-gonic/gin"
)

func SetupMarketingRouter(r *gin.Engine) {
	marketing := r.Group("/marketing")
	{
		marketing.GET("/get", marketingController.GetMarketing)
		marketing.POST("/post", marketingController.CreateMarketing)
	}
}
