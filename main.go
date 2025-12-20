// 生成一个使用gin框架的简单web服务器
package main

import (
	"project-manager/initialize"
	"project-manager/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	err := initialize.InitDB()
	if err != nil {
		panic(err)
	}
	//初始化路由
	r := routes.InitRoutes()
	//运行服务器
	runServer(r)
}

func runServer(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "欢迎使用项目管理器 服务运行正常"})
	})
	// r.Run("127.0.0.1:8086")
	if err := r.Run("127.0.0.1:8086"); err != nil {
		panic(err)
	}
}
