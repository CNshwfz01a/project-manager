package routes

import (
	"project-manager/controllers"
	//使用AuthRequired中间件
	"project-manager/middleware"

	"github.com/gin-gonic/gin"
)

// role相关请求的路由
func InitRoleRoutes(r *gin.RouterGroup) gin.IRoutes {
	role := r.Group("/roles")
	// 使用AuthRequired中间件验证用户会话
	role.Use(middleware.AuthRequired())
	{
		role.GET("/", controllers.Role.List)
		role.POST("/", controllers.Role.Add)
		// r.PUT("/roles/:id", controllers.Role.Update)
		role.DELETE("/:id", controllers.Role.Delete)
	}

	return r
}
