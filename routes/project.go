package routes

import (
	projectController "ayana/controllers/project"

	"github.com/gin-gonic/gin"
)

func SetupProjectRouter(r *gin.Engine) {
	project := r.Group("/project")
	{
		project.GET("/get", projectController.GetProject)
		project.GET("/getById/:id", projectController.GetProjectById)
		project.POST("/post", projectController.CreateProject)
		project.PUT("/project-finish/:id", projectController.FinishProject)
		project.PUT("/edit", projectController.EditProject)
		project.DELETE("/delete/:id", projectController.DeleteProject)
	}
}
