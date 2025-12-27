package service

import (
	"fmt"
	"project-manager/model"
	"project-manager/model/request"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
)

type AuditService struct{}

// List 查询审计日志列表
func (s *AuditService) List(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.AuditListReq)
	if !ok {
		return nil, ReqAssertErr
	}

	// 权限检查：仅 admin 可以查询审计日志
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		return nil, pkg.NewUnauthorizedError()
	}

	// 查询审计日志
	audits, total, err := model.AuditData.List(r.Keyword, r.StartAt, r.EndAt, r.OrderBy, r.Page, r.PageSize)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("查询审计日志失败: %v", err))
	}

	// 格式化返回
	return map[string]any{
		"list":  audits,
		"total": total,
	}, nil
}
