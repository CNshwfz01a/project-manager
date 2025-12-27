package controllers

import (
	"project-manager/model/request"
	"project-manager/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectController struct{}

// AddUserToProject 添加用户到项目
func (m *ProjectController) AddUserToProject(c *gin.Context) {
	req := new(request.ProjectAddUserReq)
	//读url的project_id参数
	var id, _ = strconv.Atoi(c.Param("project_id"))
	req.ProjectID = uint(id)
	Handle(c, req, func() (any, any) {
		return service.Project.AddUserToProject(c, req)
	})
}

// UpdateProject 更新项目信息
func (m *ProjectController) UpdateProject(c *gin.Context) {
	req := new(request.ProjectUpdateReq)
	//读url的project_id参数
	var id, _ = strconv.Atoi(c.Param("project_id"))
	Handle(c, req, func() (any, any) {
		return service.Project.UpdateProject(c, id, req)
	})
}
