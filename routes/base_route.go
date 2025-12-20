package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "路由不存在"})
	})
	//初始化各个模块的路由
	path := r.Group("/api")
	InitRoleRoutes(path)
	return r
}
