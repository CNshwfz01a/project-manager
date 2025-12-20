package routes

import (
	"project-manager/controllers"
	//使用AuthRequired中间件
	"github.com/gin-gonic/gin"
)

// role相关请求的路由
func InitRoleRoutes(r *gin.RouterGroup) gin.IRoutes {
	// role := r.Group("/role")
	// 使用AuthRequired中间件验证用户会话
	// r.Use(middleware.AuthRequired())
	// {
	// 	r.GET("/roles", controllers.Role.List)
	// 	r.POST("/role", controllers.Role.Add)
	// 	r.PUT("/role/:id", controllers.Role.Update)
	// 	r.DELETE("/role/:id", controllers.Role.Delete)
	// }
	r.GET("/roles", controllers.Role.List)
	r.POST("/role", controllers.Role.Add)
	r.PUT("/role/:id", controllers.Role.Update)
	r.DELETE("/role/:id", controllers.Role.Delete)
	return r
}
