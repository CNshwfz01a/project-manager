package routes

import (
	//使用AuthRequired中间件
	"project-manager/controllers"
	"project-manager/middleware"

	"github.com/gin-gonic/gin"
)

// user相关
func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	//登录
	r.POST("/login", controllers.User.Login)
	//登出
	r.POST("/logout", controllers.User.Logout).Use(middleware.AuthRequired())
	user := r.Group("/users")
	//使用AuthRequired中间件验证用户会话
	user.Use(middleware.AuthRequired())
	{
		user.POST("/", controllers.User.Add)

		//分配用户角色
		user.POST("/:id/roles", controllers.User.AssignRole)
		//移除用户角色
		user.DELETE("/:id/roles/:role_id", controllers.User.RemoveRole)
		user.DELETE("/:id", controllers.User.Delete)
		user.GET("/", controllers.User.List)
		user.GET("/:id", controllers.User.Detail)
		//用户团队列表
		user.GET("/:id/teams", controllers.User.TeamList)
		//用户项目列表
		user.GET("/:id/projects", controllers.User.ProjectList)

	}
	me := r.Group("/me")
	me.Use(middleware.AuthRequired())
	{
		me.PUT("/password", controllers.User.ChangePassword)
		me.PUT("/", controllers.User.UpdateProfile)
		//查询自身信息
		me.GET("/", controllers.User.MyDetail)
		//我所在的团队列表
		me.GET("/teams", controllers.User.MyTeamList)
		//我所在的项目列表
		me.GET("/projects", controllers.User.MyProjectList)
		//退出team
		me.DELETE("/teams/:team_id", controllers.User.LeaveTeam)
		//退出project
		me.DELETE("/projects/:project_id", controllers.User.LeaveProject)
	}
	return r
}
