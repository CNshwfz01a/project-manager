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

// PartialUpdateProject 部分更新项目信息
func (m *ProjectController) PartialUpdateProject(c *gin.Context) {
	req := new(request.ProjectPatch)
	//读url的project_id参数
	var id, _ = strconv.Atoi(c.Param("project_id"))
	Handle(c, req, func() (any, any) {
		return service.Project.PartialUpdateProject(c, id, req)
	})
}

// DeleteProject 删除项目
func (m *ProjectController) DeleteProject(c *gin.Context) {
	//读url的project_id参数
	var id, _ = strconv.Atoi(c.Param("project_id"))
	Handle(c, nil, func() (any, any) {
		return service.Project.DeleteProject(c, id)
	})
}

// RemoveUserFromProject 清退项目中的用户
func (m *ProjectController) RemoveUserFromProject(c *gin.Context) {
	//读url的project_id参数
	var projectID, _ = strconv.Atoi(c.Param("project_id"))
	var userID, _ = strconv.Atoi(c.Param("user_id"))
	Handle(c, nil, func() (any, any) {
		return service.Project.RemoveUserFromProject(c, projectID, userID)
	})
}

// GetProjectUsers 获取项目中的所有用户
func (m *ProjectController) GetProjectUsers(c *gin.Context) {
	//复用UserListReq
	req := new(request.UserListReq)
	//读url的project_id参数
	var id, _ = strconv.Atoi(c.Param("project_id"))
	Handle(c, req, func() (any, any) {
		return service.Project.GetProjectUsers(c, id, req)
	})
}

// GetProjectDetail 获取项目详情
func (m *ProjectController) GetProjectDetail(c *gin.Context) {
	//读url的project_id参数
	var id, _ = strconv.Atoi(c.Param("project_id"))
	Handle(c, nil, func() (any, any) {
		return service.Project.GetProjectDetail(c, id)
	})
}
