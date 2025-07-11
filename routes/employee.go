package routes

import (
	employeeController "ayana/controllers/hr"
	"ayana/utils/handler"

	"github.com/gin-gonic/gin"
)

func SetupEmployeeRouter(r *gin.Engine) {
	employee := r.Group("/employee")
	{
		employee.GET("/get", employeeController.GetEmployees)
		employee.POST("/post", employeeController.CreateEmployee)
		employee.DELETE("delete/:id", employeeController.DeleteEmployee)
		employee.PUT("/edit/:id", employeeController.UpdateEmployee)

		employee.GET("/get/presence", employeeController.GetPresence)
		employee.POST("/upload-presence", handler.UploadPresenceHandler)

		employee.POST("/post-presence-rule", employeeController.CreatePresenceRules)
		employee.GET("/get/presence-rule", employeeController.GetPresenceRules)
		employee.PUT("/edit/presence-rule/:id", employeeController.UpdatePresenceRule)

		employee.DELETE("/delete/presence-rule/:id", employeeController.DeletePresenceRule)

	}
}
