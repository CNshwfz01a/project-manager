package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired 是一个用于验证用户会话的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		//根据userId获取用户角色
		userRole := session.Get("user_role")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权,请先登录",
			})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中,供后续处理器使用
		c.Set("user_id", userID)
		c.Set("user_role", userRole)
		c.Next()
	}
}
