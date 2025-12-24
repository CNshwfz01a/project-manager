package middleware

import (
	"github.com/gin-gonic/gin"
)

// 记录请求日志的后置中间件
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		//处理请求
		c.Next()

		//记录日志
		// log.Printf("请求方法: %s, 请求路径: %s, 状态码: %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}
