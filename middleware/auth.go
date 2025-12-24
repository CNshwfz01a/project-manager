package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"project-manager/model"
)

// AuthRequired 是一个用于验证用户会话的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取user结构体
		data, err := c.Cookie("session-login")
		// log.Printf("获取到的cookie: %s", data)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未登录或会话已过期",
			})
			c.Abort()
			return
		}
		session := sessions.Default(c)
		userID := session.Get(data)
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未登录或会话已过期",
			})
			c.Abort()
			return
		}
		user, err := model.UserData.GetByID(userID.(uint))
		//如果当前密码是初始密码并且访问的不是修改密码的接口,则拒绝访问
		if user.Status == 0 && c.FullPath() != "/api/me/password" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "请先修改初始密码",
			})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中,供后续处理器使用
		c.Set("user_id", userID)
		c.Set("user_role", user.Roles)
		c.Next()
	}
}
