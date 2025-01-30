package routes

import (
	weeklyProgressController "ayana/controllers/weeklyProgress"
	// "ayana/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupWeeklyProgressRouter(r *gin.Engine) {
	weeklyProgress := r.Group("/weeklyprogress")
	{
		weeklyProgress.GET("/getById/:id", weeklyProgressController.GetWeeklyProgressById)
		weeklyProgress.POST("/post", weeklyProgressController.CreateWeeklyProgress)
		weeklyProgress.PUT("/edit", weeklyProgressController.EditWeeklyProgress)
		// weeklyProgress.GET("getById/:id", weeklyProgress.HomeById)
		// weeklyProgress.DELETE("deleteById/:id", weeklyProgress.DeleteHome)
		// weeklyProgress.PUT("update/:id", weeklyProgress.UpdateHome)
		// weeklyProgress.POST("/img", middlewares.FileUploadMiddleware(), weeklyProgress.AddImage)
	}
}
