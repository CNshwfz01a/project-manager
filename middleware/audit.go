package middleware

import (
	"fmt"
	"project-manager/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 记录请求日志的后置中间件
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 判断是否需要记录审计日志
		if shouldAudit(c) {
			// 获取用户信息
			userID := c.GetUint("user_id")
			username := c.GetString("username")
			if username == "" {
				username = "unknown"
			}

			// 构建审计内容
			content := buildAuditContent(c, userID, username, startTime)

			// 记录到数据库
			_ = model.AuditData.Create(content)
		}
	}
}

// shouldAudit 判断是否需要记录审计日志
func shouldAudit(c *gin.Context) bool {
	method := c.Request.Method
	path := c.Request.URL.Path
	statusCode := c.Writer.Status()

	// 只记录成功的请求（2xx）
	if statusCode < 200 || statusCode >= 300 {
		return false
	}

	// 不记录查询操作（GET）
	if method == "GET" {
		return false
	}

	// 不记录健康检查
	if path == "/healthz" {
		return false
	}

	// 记录所有其他的增删改操作（POST、PUT、PATCH、DELETE）
	return method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE"
}

// buildAuditContent 构建审计内容
func buildAuditContent(c *gin.Context, userID uint, username string, startTime time.Time) string {
	method := c.Request.Method
	path := c.Request.URL.Path
	statusCode := c.Writer.Status()
	duration := time.Since(startTime).Milliseconds()
	clientIP := c.ClientIP()

	// 构建操作描述
	action := getActionDescription(method, path)

	// 格式化审计内容
	content := fmt.Sprintf(
		"[%s] 用户: %s (ID:%d) | 操作: %s | 路径: %s %s | 状态: %d | IP: %s | 耗时: %dms",
		time.Now().Format("2006-01-02 15:04:05"),
		username,
		userID,
		action,
		method,
		path,
		statusCode,
		clientIP,
		duration,
	)

	return content
}

// getActionDescription 根据请求方法和路径获取操作描述
func getActionDescription(method, path string) string {
	switch {
	// 登录登出
	case strings.Contains(path, "/api/login"):
		return "用户登录"
	case strings.Contains(path, "/api/logout"):
		return "用户登出"

	// 用户相关
	case strings.Contains(path, "/api/me/password"):
		return "修改密码"
	case strings.Contains(path, "/api/me") && method == "PUT":
		return "修改个人信息"
	case strings.Contains(path, "/api/users") && method == "POST":
		return "创建用户"
	case strings.Contains(path, "/api/users") && method == "DELETE":
		return "删除用户"
	case strings.Contains(path, "/api/users") && strings.Contains(path, "/roles") && method == "POST":
		return "添加用户角色"
	case strings.Contains(path, "/api/users") && strings.Contains(path, "/roles") && method == "DELETE":
		return "移除用户角色"

	// Team 相关
	case strings.Contains(path, "/api/teams") && method == "POST" && !strings.Contains(path, "/users") && !strings.Contains(path, "/projects"):
		return "创建团队"
	case strings.Contains(path, "/api/teams") && method == "PUT":
		return "更新团队"
	case strings.Contains(path, "/api/teams") && method == "PATCH":
		return "修改团队Leader"
	case strings.Contains(path, "/api/teams") && method == "DELETE" && !strings.Contains(path, "/users"):
		return "删除团队"
	case strings.Contains(path, "/api/teams") && strings.Contains(path, "/users") && method == "POST":
		return "添加团队成员"
	case strings.Contains(path, "/api/teams") && strings.Contains(path, "/users") && method == "DELETE":
		return "移除团队成员"
	case strings.Contains(path, "/api/teams") && strings.Contains(path, "/projects") && method == "POST":
		return "创建项目"
	case strings.Contains(path, "/api/me/teams") && method == "DELETE":
		return "退出团队"

	// Project 相关
	case strings.Contains(path, "/api/projects") && method == "PUT":
		return "更新项目"
	case strings.Contains(path, "/api/projects") && method == "PATCH":
		return "部分更新项目"
	case strings.Contains(path, "/api/projects") && method == "DELETE" && !strings.Contains(path, "/users"):
		return "删除项目"
	case strings.Contains(path, "/api/projects") && strings.Contains(path, "/users") && method == "POST":
		return "添加项目成员"
	case strings.Contains(path, "/api/projects") && strings.Contains(path, "/users") && method == "DELETE":
		return "移除项目成员"
	case strings.Contains(path, "/api/me/projects") && method == "DELETE":
		return "退出项目"

	// Role 相关
	case strings.Contains(path, "/api/roles") && method == "POST":
		return "创建角色"
	case strings.Contains(path, "/api/roles") && method == "DELETE":
		return "删除角色"

	default:
		return fmt.Sprintf("%s操作", method)
	}
}
