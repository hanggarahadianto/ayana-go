package routes

import (
	customerController "ayana/controllers/customer"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRouter(r *gin.Engine) {
	customer := r.Group("/customer")
	{
		customer.GET("/get", customerController.GetCustomers)
		customer.POST("/post", customerController.CreateCustomer)
		customer.PUT("/update/:id", customerController.UpdateCustomer)

		// cluster.PUT("/update/:id", clusterController.UpdateCluster)
		// cluster.GET("/getById/:id", clusterController.GetClusterByID)
		// cluster.DELETE("deleteById/:id", clusterController.DeleteCluster)

		// home.POST("/:homeId/images", homeController.UploadProductImage)
		// home.GET("/:homeId/images", homeController.GetHomeImages)
	}
}
