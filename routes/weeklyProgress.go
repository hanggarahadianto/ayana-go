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
	}
}
