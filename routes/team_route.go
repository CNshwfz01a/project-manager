package routes

import (
	"project-manager/controllers"
	"project-manager/middleware"

	"github.com/gin-gonic/gin"
)

func InitTeamRoutes(r *gin.RouterGroup) gin.IRoutes {
	team := r.Group("/teams")
	team.Use(middleware.AuthRequired())
	{
		team.POST("/", controllers.Team.Add)
		team.POST("/:team_id/users", controllers.Team.AddUserToTeam)       //加人
		team.POST("/:team_id/projects", controllers.Team.AddProjectToTeam) //加项目
		//修改单一属性
		team.PATCH("/:team_id", controllers.Team.Patch)
		//修改属性
		team.PUT("/:team_id", controllers.Team.Update)
		//查询详情
		team.GET("/:team_id", controllers.Team.Get)
		//查询列表
		team.GET("/", controllers.Team.List)
		//查询团队用户列表
		team.GET("/:team_id/users", controllers.Team.ListUsers)
		//查询团队项目列表
		team.GET("/:team_id/projects", controllers.Team.ListProjects)
		//删除
		team.DELETE("/:team_id", controllers.Team.Delete)
		//从团队移除用户
		team.DELETE("/:team_id/users/:user_id", controllers.Team.RemoveUserFromTeam)
	}
	return r
}
