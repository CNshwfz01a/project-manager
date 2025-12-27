package controllers

import (
	"project-manager/model/request"
	"project-manager/service"

	"github.com/gin-gonic/gin"
)

type AuditController struct{}

// List 查询审计日志列表
func (m *AuditController) List(c *gin.Context) {
	req := new(request.AuditListReq)
	req.SetDefaults()
	req.SetParams(c)

	Handle(c, req, func() (any, any) {
		return service.Audit.List(c, req)
	})
}
