package routes

import (
	companyController "ayana/controllers/company"
	middlewares "ayana/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func SetupCompanyRouter(r *gin.Engine) {
	company := r.Group("/company", middlewares.AuthMiddleware())
	{
		company.GET("/get", companyController.GetCompany)
		company.GET("/get-by-user", companyController.GetCompaniesByUser)
		company.POST("/post", companyController.CreateCompany)
		company.PUT("/update/:id", companyController.UpdateCompany)
		company.DELETE("/delete/:id", companyController.DeleteCompany)
		company.POST("/post/assign-user", companyController.AssignCompanyToUsers)
		// project.PUT("/edit", projectController.EditProject)
		// project.DELETE("/delete/:id", projectController.DeleteProject)
	}
}
