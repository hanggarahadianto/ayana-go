package routes

import (
	homeController "ayana/controllers/home"
	"ayana/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupHomeRouter(r *gin.Engine) {
	home := r.Group("/home")
	{
		home.GET("/get", homeController.GetHomes)
		home.POST("/post", middlewares.FileUploadMiddleware(), homeController.CreateHome)
		home.GET("getById/:id", homeController.HomeById)
		home.DELETE("deleteById/:id", homeController.DeleteHome)
		home.PUT("update/:id", homeController.UpdateHome)
		home.POST("/img", middlewares.FileUploadMiddleware(), homeController.AddImage)
	}
}
