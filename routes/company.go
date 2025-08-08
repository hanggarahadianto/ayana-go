package routes

import (
	companyController "ayana/controllers/company"

	"github.com/gin-gonic/gin"
)

func SetupCompanyRouter(r *gin.Engine) {
	company := r.Group("/company")
	{
		company.GET("/get", companyController.GetCompany)
		company.GET("/get-by-user", companyController.GetCompaniesByUser)
		company.POST("/post", companyController.CreateCompany)
		company.POST("/post/assign-user", companyController.AssignCompanyToUsers)
		// project.PUT("/edit", projectController.EditProject)
		// project.DELETE("/delete/:id", projectController.DeleteProject)
	}
}
