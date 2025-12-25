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
		team.POST("/:team_id/users", controllers.Team.AddUserToTeam)
		team.POST("/:team_id/projects", controllers.Team.AddProjectToTeam)
		//修改单一属性
		team.PATCH("/:team_id", controllers.Team.Patch)
	}
	return r
}
