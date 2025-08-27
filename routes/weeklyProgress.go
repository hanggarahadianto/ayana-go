package routes

import (
	weeklyProgressController "ayana/controllers/weeklyProgress"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupWeeklyProgressRouter(r *gin.Engine) {
	weeklyProgress := r.Group("/weeklyprogress", middlewares.AuthMiddleware())
	{
		weeklyProgress.GET("/getById/:id", weeklyProgressController.GetWeeklyProgressById)
		weeklyProgress.POST("/post", weeklyProgressController.CreateWeeklyProgress)
		weeklyProgress.PUT("/edit", weeklyProgressController.EditWeeklyProgress)
		weeklyProgress.DELETE("delete/:id", weeklyProgressController.DeleteWeeklyProgress)
		//
		// weeklyProgress.GET("getById/:id", weeklyProgress.HomeById)
		// weeklyProgress.PUT("update/:id", weeklyProgress.UpdateHome)
		// weeklyProgress.POST("/img", middlewares.FileUploadMiddleware(), weeklyProgress.AddImage)
	}
}
