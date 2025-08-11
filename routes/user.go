package routes

import (
	userController "ayana/controllers/user" // alias biar jelas ini controller
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.Engine) {
	userRoutes := r.Group("/user", middlewares.AuthMiddleware())
	{
		userRoutes.GET("/get-by-id/:id", userController.GetUsers)
		userRoutes.DELETE("delete/:id", userController.DeleteUser)
		userRoutes.PUT("update/:id", userController.UpdateUser)

	}
}
