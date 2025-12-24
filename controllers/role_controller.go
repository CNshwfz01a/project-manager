package controllers

import (
	"project-manager/model/request"
	"project-manager/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// GetRoles 获取所有角色列表
func (m *RoleController) List(c *gin.Context) {
	req := new(request.RoleListReq)
	Handle(c, req, func() (any, any) {
		return service.Role.List(c, req)
	})
}

// AddRole 添加新角色 只有管理员可以添加角色 并且角色名不能重复
func (m *RoleController) Add(c *gin.Context) {
	req := new(request.RoleAddReq)
	Handle(c, req, func() (any, any) {
		return service.Role.Add(c, req)
	})
}

// UpdateRole 更新角色信息
// func (m *RoleController) Update(c *gin.Context) {
// }

// DeleteRole 删除角色
func (m *RoleController) Delete(c *gin.Context) {
	req := new(request.RoleDeleteReq)
	var id, _ = strconv.Atoi(c.Param("id"))
	req.ID = uint(id)
	Handle(c, req, func() (any, any) {
		return service.Role.Delete(c, req)
	})
}
