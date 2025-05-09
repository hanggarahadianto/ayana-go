package routes

import (
	homeController "ayana/controllers/home"

	"github.com/gin-gonic/gin"
)

func SetupHomeRouter(r *gin.Engine) {
	home := r.Group("/home")
	{
		home.GET("/get", homeController.GetHomes)
		home.POST("/post", homeController.CreateHome)
		home.GET("getById/:id", homeController.HomeById)
		home.GET("/getByClusterId/:cluster_id", homeController.HomeListByClusterId) //
		home.DELETE("deleteById/:id", homeController.DeleteHome)

		home.POST("/:homeId/images", homeController.UploadProductImage)
		home.GET("/:homeId/images", homeController.GetHomeImages)
	}
}
