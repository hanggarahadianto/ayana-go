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
		cluster.PUT("/update/:id", clusterController.UpdateCluster)
		cluster.GET("/getById/:id", clusterController.GetClusterByID)
		cluster.DELETE("deleteById/:id", clusterController.DeleteCluster)

		// home.POST("/:homeId/images", homeController.UploadProductImage)
		// home.GET("/:homeId/images", homeController.GetHomeImages)
	}
}
