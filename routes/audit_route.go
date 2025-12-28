package routes

import (
	"project-manager/controllers"

	"github.com/gin-gonic/gin"

	"project-manager/middleware"
)

func InitAuditRoutes(r *gin.Engine) {
	auditGroup := r.Group("/api/audits")
	{
		//使用鉴权中间件
		auditGroup.GET("", controllers.Audit.List).Use(middleware.AuthRequired())
	}
}
