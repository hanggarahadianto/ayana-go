package routes

import (
	payoutController "ayana/controllers/payout"

	"github.com/gin-gonic/gin"
)

func SetupPayoutRouter(r *gin.Engine) {
	payout := r.Group("/payout")
	{
		payout.GET("/get", payoutController.GetPayoutsByCompany)
		// project.GET("/getById/:id", projectController.GetProjectById)
		payout.POST("/post", payoutController.CreatePayout)
		payout.PUT("/edit", payoutController.EditPayout)
		payout.DELETE("/delete/:id", payoutController.DeletePayout)
	}
}
