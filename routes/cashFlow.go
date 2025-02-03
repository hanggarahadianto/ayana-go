package routes

import (
	cashFlowController "ayana/controllers/cashFlow"
	// "ayana/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCashFlowRouter(r *gin.Engine) {
	cashFlow := r.Group("/cashflow")
	{
		cashFlow.GET("/getById/:id", cashFlowController.GetCashFlowById)
		cashFlow.POST("/post", cashFlowController.CreateCashFlow)
		// weeklyProgress.PUT("/edit", cashFlowController.EditW)
		// weeklyProgress.GET("getById/:id", weeklyProgress.HomeById)
		// weeklyProgress.DELETE("deleteById/:id", weeklyProgress.DeleteHome)
		// weeklyProgress.PUT("update/:id", weeklyProgress.UpdateHome)
		// weeklyProgress.POST("/img", middlewares.FileUploadMiddleware(), weeklyProgress.AddImage)
	}
}
