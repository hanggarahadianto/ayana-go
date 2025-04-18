package routes

import (
	payinController "ayana/controllers/payin"

	"github.com/gin-gonic/gin"
)

func SetupPayinRouter(r *gin.Engine) {
	payin := r.Group("/payin")
	{
		payin.POST("/post", payinController.CreatePayIn)

	}
}
