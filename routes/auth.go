package routes

import (
	authController "ayana/controllers/auth"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/logout", authController.Logout)
	}
}
