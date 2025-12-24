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
		user.GET("/", controllers.User.List)
		// user.GET("/users/:id", controllers.User.Detail)

		// user.PUT("/users/:id", controllers.User.Update)
		// user.DELETE("/users/:id", controllers.User.Delete)
	}
	me := r.Group("/me")
	me.Use(middleware.AuthRequired())
	{
		me.PUT("/password", controllers.User.ChangePassword)
	}
	return r
}
