package routes

import (
	companyController "ayana/controllers/company"

	"github.com/gin-gonic/gin"
)

func SetupCompanyRouter(r *gin.Engine) {
	company := r.Group("/company")
	{
		company.GET("/get", companyController.GetCompany)
		// project.GET("/getById/:id", projectController.GetProjectById)
		company.POST("/post", companyController.CreateCompany)
		// project.PUT("/edit", projectController.EditProject)
		// project.DELETE("/delete/:id", projectController.DeleteProject)
	}
}
