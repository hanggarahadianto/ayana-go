package routes

import (
	clusterController "ayana/controllers/cluster"

	"github.com/gin-gonic/gin"
)

func SetupClusterRouter(r *gin.Engine) {
	cluster := r.Group("/cluster")
	{
		cluster.GET("/get", clusterController.GetCluster)
		cluster.POST("/post", clusterController.CreateCluster)
		cluster.GET("/getById/:id", clusterController.GetClusterByID)
		// home.DELETE("deleteById/:id", homeController.DeleteHome)

		// home.POST("/:homeId/images", homeController.UploadProductImage)
		// home.GET("/:homeId/images", homeController.GetHomeImages)
	}
}
