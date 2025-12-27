package routes

import (
	"project-manager/controllers"

	"github.com/gin-gonic/gin"
)

func InitAuditRoutes(r *gin.Engine) {
	auditGroup := r.Group("/api/audits")
	{
		auditGroup.GET("", controllers.Audit.List)
	}
}
