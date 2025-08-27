package routes

import (
	cashFlowController "ayana/controllers/cashFlow"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupCashFlowRouter(r *gin.Engine) {
	cashFlow := r.Group("/cashflow", middlewares.AuthMiddleware())
	{
		cashFlow.GET("/getByProjectId/:id", cashFlowController.GetCashFlowListByProjectId)
		cashFlow.POST("/post", cashFlowController.CreateCashFlow)
		cashFlow.PUT("/edit/:id", cashFlowController.UpdateCashFlow)
		// weeklyProgress.GET("getById/:id", weeklyProgress.HomeById)
		cashFlow.DELETE("deleteById/:id", cashFlowController.DeleteCashFlow)
		// weeklyProgress.PUT("update/:id", weeklyProgress.UpdateHome)
		// weeklyProgress.POST("/img", middlewares.FileUploadMiddleware(), weeklyProgress.AddImage)
	}
}
