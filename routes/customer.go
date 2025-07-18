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
		customer.DELETE("/delete/:id", customerController.DeleteCustomer)

		customer.GET("testimony/get", customerController.GetAllTestimonies)
		customer.POST("testimony/post", customerController.CreateCustomerTestimony)
		customer.PUT("testimony/update/:id", customerController.UpdateCustomerTestimony)
		customer.DELETE("testimony/delete/:id", customerController.DeleteTestimony)

	}
}
