package routes

import (
	goodController "ayana/controllers/good"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupGoodRouter(r *gin.Engine) {
	good := r.Group("/good", middlewares.AuthMiddleware())
	{
		good.GET("/getByCashFlowId", goodController.GetGoodByCashFlowId)
		good.POST("/post", goodController.CreateGood)
		good.PUT("/edit", goodController.UpdateGood)
		// weeklyProgress.GET("getById/:id", weeklyProgress.HomeById)

		// weeklyProgress.PUT("update/:id", weeklyProgress.UpdateHome)
		// weeklyProgress.POST("/img", middlewares.FileUploadMiddleware(), weeklyProgress.AddImage)
	}
}
