package routes

import (
	"project-manager/controllers"
	"project-manager/middleware"

	"github.com/gin-gonic/gin"
)

func InitProjectRoute(r *gin.RouterGroup) gin.IRoutes {
	project := r.Group("/projects")
	project.Use(middleware.AuthRequired())
	{
		project.POST("/:project_id/users", controllers.Project.AddUserToProject)
	}
	return r
}
