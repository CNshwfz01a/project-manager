package routes

import (
	"project-manager/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()
	//使用cookie存储session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("project_manager_session", store))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "路由不存在"})
	})
	//使用audit中间件记录用户操作日志
	r.Use(middleware.Audit())
	//初始化各个模块的路由
	path := r.Group("/api")
	InitUserRoutes(path)
	InitRoleRoutes(path)
	InitTeamRoutes(path)
	InitProjectRoute(path)
	InitAuditRoutes(path)
	return r
}
