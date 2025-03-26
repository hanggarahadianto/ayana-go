package routes

import (
	cashFlowController "ayana/controllers/cashFlow"
	// "ayana/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCashFlowRouter(r *gin.Engine) {
	cashFlow := r.Group("/cashflow")
	{
		cashFlow.GET("/getByProjectId/:id", cashFlowController.GetCashFlowListByProjectId)
		cashFlow.POST("/post", cashFlowController.CreateCashFlow)
		cashFlow.PUT("/edit/:id", cashFlowController.UpdateCashFlow)
		// weeklyProgress.GET("getById/:id", weeklyProgress.HomeById)
		// weeklyProgress.DELETE("deleteById/:id", weeklyProgress.DeleteHome)
		// weeklyProgress.PUT("update/:id", weeklyProgress.UpdateHome)
		// weeklyProgress.POST("/img", middlewares.FileUploadMiddleware(), weeklyProgress.AddImage)
	}
}
