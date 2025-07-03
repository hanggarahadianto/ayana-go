package routes

import (
	employeeController "ayana/controllers/hr"

	"github.com/gin-gonic/gin"
)

func SetupEmployeeRouter(r *gin.Engine) {
	employee := r.Group("/employee")
	{
		employee.GET("/get", employeeController.GetEmployees)
		employee.POST("/post", employeeController.CreateEmployee)
		employee.DELETE("delete/:id", employeeController.DeleteEmployee)
		employee.PUT("/edit/:id", employeeController.UpdateEmployee)

	}
}
