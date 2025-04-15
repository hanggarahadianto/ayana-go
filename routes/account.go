package routes

import (
	accountController "ayana/controllers/account"
	// "ayana/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAccountRouter(r *gin.Engine) {
	account := r.Group("/account")
	{
		account.GET("/get", accountController.GetAccount)
		account.POST("/post", accountController.CreateAccount)

	}
}
