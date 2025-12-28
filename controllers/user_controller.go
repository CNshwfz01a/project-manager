package controllers

import (
	"project-manager/model/request"
	"project-manager/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (m *UserController) Login(c *gin.Context) {
	req := new(request.UserLoginReq)
	Handle(c, req, func() (any, any) {
		return service.User.Login(c, req)
	})
}

func (m *UserController) Logout(c *gin.Context) {
	Handle(c, nil, func() (any, any) {
		return service.User.Logout(c)
	})
}

// 修改自身密码 修改完成需要重新登录
func (m *UserController) ChangePassword(c *gin.Context) {
	req := new(request.UserChangePasswordReq)
	Handle(c, req, func() (any, any) {
		return service.User.ChangePassword(c, req)
	})
}

func (m *UserController) Add(c *gin.Context) {
	req := new(request.UserAddReq)
	Handle(c, req, func() (any, any) {
		return service.User.Add(c, req)
	})
}

// 查询用户列表
func (m *UserController) List(c *gin.Context) {
	req := new(request.UserListReq)
	req.SetParams(c)
	Handle(c, req, func() (any, any) {
		return service.User.List(c, req)
	})
}

// 删除用户
func (m *UserController) Delete(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	Handle(c, nil, func() (any, any) {
		return service.User.Delete(c, uint(id))
	})
}

// 分配用户角色
func (m *UserController) AssignRole(c *gin.Context) {
	req := new(request.UserAssignRoleReq)
	var id, _ = strconv.Atoi(c.Param("id"))
	Handle(c, req, func() (any, any) {
		return service.User.AssignRole(c, req, uint(id))
	})
}

// 移除用户角色
func (m *UserController) RemoveRole(c *gin.Context) {
	var userID, _ = strconv.Atoi(c.Param("id"))
	var roleID, _ = strconv.Atoi(c.Param("role_id"))
	Handle(c, nil, func() (any, any) {
		return service.User.RemoveRole(c, uint(userID), uint(roleID))
	})
}

// 查询自身信息
func (m *UserController) MyDetail(c *gin.Context) {
	Handle(c, nil, func() (any, any) {
		return service.User.MyDetail(c)
	})
}

// 更新自身资料
func (m *UserController) UpdateProfile(c *gin.Context) {
	req := new(request.UserUpdateProfileReq)
	Handle(c, req, func() (any, any) {
		return service.User.UpdateProfile(c, req)
	})
}

// 查询用户详情
func (m *UserController) Detail(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	Handle(c, nil, func() (any, any) {
		return service.User.Detail(c, uint(id))
	})
}

// 用户团队列表
func (m *UserController) TeamList(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	req := new(request.UserListReq)
	req.SetParams(c)
	Handle(c, req, func() (any, any) {
		return service.User.TeamList(c, uint(id), req)
	})
}

// 用户项目列表
func (m *UserController) ProjectList(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	Handle(c, nil, func() (any, any) {
		return service.User.ProjectList(c, uint(id))
	})
}

// MyTeamList 我所在的团队列表
func (m *UserController) MyTeamList(c *gin.Context) {
	//获取参数leading bool
	leadingStr := c.Query("leading")
	//为空返回0 true返回1 false返回-1
	leading := 0
	switch leadingStr {
	case "true":
		leading = 1
	case "false":
		leading = -1
	}
	Handle(c, nil, func() (any, any) {
		return service.User.MyTeamList(c, leading)
	})
}

// MyProjectList 我所在的项目列表
func (m *UserController) MyProjectList(c *gin.Context) {
	Handle(c, nil, func() (any, any) {
		return service.User.MyProjectList(c)
	})
}

// LeaveTeam 退出团队
func (m *UserController) LeaveTeam(c *gin.Context) {
	var teamID, _ = strconv.Atoi(c.Param("team_id"))
	Handle(c, nil, func() (any, any) {
		return service.User.LeaveTeam(c, uint(teamID))
	})
}

// LeaveProject 退出项目
func (m *UserController) LeaveProject(c *gin.Context) {
	var projectID, _ = strconv.Atoi(c.Param("project_id"))
	Handle(c, nil, func() (any, any) {
		return service.User.LeaveProject(c, uint(projectID))
	})
}
