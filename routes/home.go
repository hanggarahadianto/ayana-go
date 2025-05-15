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
		home.DELETE("deleteById/:homeId", homeController.DeleteHome)
		home.PUT("/update", homeController.UpdateHome)
		home.POST("/create/images/:homeId", homeController.UploadProductImage)
		home.PUT("/update/images/:homeId", homeController.UpdateProductImages)
		home.GET("/:homeId/images", homeController.GetHomeImages)
	}
}
