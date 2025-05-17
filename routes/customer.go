package routes

import (
	customerController "ayana/controllers/customer"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRouter(r *gin.Engine) {
	cluster := r.Group("/customer")
	{
		cluster.GET("/get", customerController.GetCustomers)
		cluster.POST("/post", customerController.CreateCustomer)
		// cluster.PUT("/update/:id", clusterController.UpdateCluster)
		// cluster.GET("/getById/:id", clusterController.GetClusterByID)
		// cluster.DELETE("deleteById/:id", clusterController.DeleteCluster)

		// home.POST("/:homeId/images", homeController.UploadProductImage)
		// home.GET("/:homeId/images", homeController.GetHomeImages)
	}
}
