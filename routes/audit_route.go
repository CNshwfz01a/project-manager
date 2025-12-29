package routes

import (
	"project-manager/controllers"

	"github.com/gin-gonic/gin"

	"project-manager/middleware"
)

func InitAuditRoutes(r *gin.RouterGroup) gin.IRoutes {
	auditGroup := r.Group("/audits")
	{
		//使用鉴权中间件
		auditGroup.GET("", controllers.Audit.List).Use(middleware.AuthRequired())
	}
	return r
}
