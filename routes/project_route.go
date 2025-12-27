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
		//更新项目信息
		project.PUT("/:project_id", controllers.Project.UpdateProject)
		// //更新部分项目信息
		project.PATCH("/:project_id", controllers.Project.PartialUpdateProject)
		// //删除项目
		project.DELETE("/:project_id", controllers.Project.DeleteProject)
		//清退项目中的用户
		project.DELETE("/:project_id/users/:user_id", controllers.Project.RemoveUserFromProject)
		// //获取项目详情
		project.GET("/:project_id", controllers.Project.GetProjectDetail)
		// //获取项目中的所有用户
		project.GET("/:project_id/users", controllers.Project.GetProjectUsers)
	}
	return r
}
