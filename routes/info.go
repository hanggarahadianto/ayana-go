package routes

import (
	infoController "ayana/controllers/info"

	"github.com/gin-gonic/gin"
)

func SetupInfoRouter(r *gin.Engine) {
	additionalInfo := r.Group("/info")
	{
		additionalInfo.GET("/get/:id", infoController.GetInfo)
		additionalInfo.POST("/post", infoController.CreateInfo)
	}
}
